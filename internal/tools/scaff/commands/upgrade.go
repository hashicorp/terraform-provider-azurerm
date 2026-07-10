package commands

import (
	"errors"
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/gen"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/list-upgrade"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
	"github.com/mitchellh/cli"
)

// UpgradeCommand upgrades an existing resource file so it can support List:
// it adds Resource Identity and/or refactors Read into a reusable flatten
// method, and optionally generates the list resource.
type UpgradeCommand struct {
	Ui cli.Ui
}

var _ cli.Command = UpgradeCommand{}

type upgradeData struct {
	File          string
	ARMType       string
	Service       string
	Resource      string
	APIVersion    string
	Name          string
	GoName        string
	ReadModel     string
	IdentityProps string
	Identity      bool
	Flatten       bool
	List          bool
	Write         bool
	Overwrite     bool
	PandoraURL    string
	Provider      string
	Org           string
	InputFile     string
}

// upgradeFile is the schema of an HCL file supplied via -input. It describes one
// or more existing resources to upgrade, plus optional global overrides applied
// to every resource in the file.
type upgradeFile struct {
	PandoraURL string         `hcl:"pandora_url,optional"`
	Provider   string         `hcl:"provider,optional"`
	Org        string         `hcl:"org,optional"`
	Write      *bool          `hcl:"write,optional"`
	Overwrite  *bool          `hcl:"overwrite,optional"`
	Resources  []upgradeBlock `hcl:"resource,block"`
}

// upgradeBlock is a single "resource" block within an upgrade input file. The
// block label is the terraform resource name (snake_case, without the provider
// prefix), e.g. resource "monitor_workspace" { ... }.
type upgradeBlock struct {
	Name          string `hcl:"name,label"`
	File          string `hcl:"file"`
	ARMType       string `hcl:"arm_type,optional"`
	Service       string `hcl:"service,optional"`
	Resource      string `hcl:"pandora_resource,optional"`
	APIVersion    string `hcl:"api_version,optional"`
	GoName        string `hcl:"go_name,optional"`
	ReadModel     string `hcl:"read_model,optional"`
	IdentityProps string `hcl:"identity_properties,optional"`
	Identity      *bool  `hcl:"identity,optional"`
	Flatten       *bool  `hcl:"flatten,optional"`
	List          *bool  `hcl:"list,optional"`
}

func (d *upgradeData) parseArgs(args []string) error {
	set := flag.NewFlagSet("upgrade", flag.ExitOnError)
	set.StringVar(&d.File, "file", "", "(Required) path to the existing resource .go file to upgrade")
	set.StringVar(&d.ARMType, "arm-type", "", "(Optional) ARM resource type, e.g. \"Microsoft.Monitor/accounts\"; used to resolve list operations and the SDK read model")
	set.StringVar(&d.Service, "service", "", "(Optional) explicit Pandora service name, e.g. \"Monitor\"")
	set.StringVar(&d.Resource, "resource", "", "(Optional) explicit Pandora resource key, e.g. \"AzureMonitorWorkspaces\"")
	set.StringVar(&d.APIVersion, "api-version", "", "(Optional) API version; defaults to the latest non-preview version")
	set.StringVar(&d.Name, "name", "", "(Optional) terraform resource name (snake_case, without provider prefix); derived from the file when omitted")
	set.StringVar(&d.GoName, "go-name", "", "(Optional) Go identifier override for the resource")
	set.StringVar(&d.ReadModel, "read-model", "", "(Optional) SDK read-model type name; overrides the value resolved from Pandora")
	set.StringVar(&d.IdentityProps, "identity-properties", "", "(Optional) -properties value for the identity test go:generate directive; defaults to \"name,resource_group_name\"")
	set.BoolVar(&d.Identity, "identity", false, "(Optional) add Resource Identity if missing")
	set.BoolVar(&d.Flatten, "flatten", false, "(Optional) refactor Read into a reusable flatten method if missing")
	set.BoolVar(&d.List, "list", false, "(Optional) make the resource list-ready (identity + flatten) and generate the list resource")
	set.BoolVar(&d.Write, "write", false, "(Optional) write the changes to disk; without it the command performs a dry run")
	set.BoolVar(&d.Overwrite, "overwrite", false, "(Optional) overwrite an existing generated list file")
	set.StringVar(&d.PandoraURL, "pandora-url", "", "(Optional) Pandora Data API base URL; defaults to the configured schema_api_url")
	set.StringVar(&d.Provider, "provider", "", "(Optional) provider name override, e.g. \"azurerm\"")
	set.StringVar(&d.Org, "org", "", "(Optional) provider GitHub org override, e.g. \"hashicorp\"")
	set.StringVar(&d.InputFile, "input", "", "(Optional) path to an HCL file describing one or more resources to upgrade; when set, the per-resource flags are ignored")

	if err := set.Parse(args); err != nil {
		return err
	}
	if d.InputFile == "" && d.File == "" {
		return errors.New("-file is required (or provide -input)")
	}
	return nil
}

func (c UpgradeCommand) Run(args []string) int {
	d := &upgradeData{}
	if err := d.parseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	var err error
	if d.InputFile != "" {
		err = c.runFromFile(d)
	} else {
		err = c.run(d)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("upgrading resource: %+v", err))
		return 1
	}
	return 0
}

// runFromFile reads an HCL input file and upgrades every resource block it
// declares. Global attributes in the file override the .scaff.hcl config
// defaults; the per-resource CLI flags are ignored in this mode.
func (c UpgradeCommand) runFromFile(d *upgradeData) error {
	var file upgradeFile
	if err := hclsimple.DecodeFile(d.InputFile, nil, &file); err != nil {
		return fmt.Errorf("reading input file %q: %w", d.InputFile, err)
	}
	if len(file.Resources) == 0 {
		return fmt.Errorf("input file %q defines no resource blocks", d.InputFile)
	}

	write := d.Write
	if file.Write != nil {
		write = *file.Write
	}
	overwrite := d.Overwrite
	if file.Overwrite != nil {
		overwrite = *file.Overwrite
	}

	for _, b := range file.Resources {
		if b.File == "" {
			return fmt.Errorf("resource %q: file is required", b.Name)
		}
		rd := &upgradeData{
			File:          b.File,
			ARMType:       b.ARMType,
			Service:       b.Service,
			Resource:      b.Resource,
			APIVersion:    b.APIVersion,
			Name:          b.Name,
			GoName:        b.GoName,
			ReadModel:     b.ReadModel,
			IdentityProps: b.IdentityProps,
			Identity:      derefBool(b.Identity),
			Flatten:       derefBool(b.Flatten),
			List:          derefBool(b.List),
			Write:         write,
			Overwrite:     overwrite,
			PandoraURL:    file.PandoraURL,
			Provider:      file.Provider,
			Org:           file.Org,
		}
		if err := c.run(rd); err != nil {
			return fmt.Errorf("resource %q: %w", b.Name, err)
		}
	}
	return nil
}

func derefBool(b *bool) bool {
	return b != nil && *b
}

func (c UpgradeCommand) run(d *upgradeData) error {
	res, err := list_upgrade.Analyze(d.File)
	if err != nil {
		return err
	}

	plan := res.PlanUpgrade()
	c.printPlan(d.File, plan)
	if !plan.Supported {
		return fmt.Errorf("%s", plan.Reason)
	}

	// With no action flags this is a read-only report.
	if !d.Identity && !d.Flatten && !d.List {
		c.Ui.Info("no action requested; pass -identity, -flatten and/or -list (add -write to apply)")
		return nil
	}

	wantIdentity := (d.Identity || d.List) && !res.HasIdentity
	wantFlatten := (d.Flatten || d.List) && !res.HasFlatten

	// Resolve the IR from Pandora when we need the read model or list operations.
	var resolved *ir.ResourceIR
	if d.List || (wantFlatten && d.ReadModel == "") {
		resolved, err = c.resolveIR(d, res)
		if err != nil {
			return fmt.Errorf("resolving API definitions (provide -read-model to skip, or -arm-type/-service+-resource): %w", err)
		}
	}

	readModel := d.ReadModel
	if readModel == "" && resolved != nil {
		readModel = resolved.ReadModel
	}
	if wantFlatten && readModel == "" {
		return errors.New("refactoring Read into flatten needs the SDK read model; provide -read-model or -arm-type")
	}

	name := d.Name
	if name == "" {
		name = strings.TrimPrefix(res.TerraformType, d.provider()+"_")
	}

	newSrc, changed, err := res.Upgrade(list_upgrade.UpgradeOptions{
		AddIdentity:        wantIdentity,
		ExtractFlatten:     wantFlatten,
		ReadModel:          readModel,
		ResourceName:       name,
		IdentityProperties: d.IdentityProps,
	})
	if err != nil {
		return err
	}

	if err := c.applyResource(d, res.Path, newSrc, changed); err != nil {
		return err
	}

	if d.List {
		if resolved == nil {
			return errors.New("-list requires API definitions; provide -arm-type or -service and -resource")
		}
		reconcileIR(resolved, res)
		if err := c.generateList(d, resolved); err != nil {
			return err
		}
	}
	return nil
}

// applyResource writes or previews the upgraded resource file.
func (c UpgradeCommand) applyResource(d *upgradeData, path string, newSrc []byte, changed bool) error {
	if !changed {
		c.Ui.Info(fmt.Sprintf("%s: already up to date", path))
		return nil
	}
	if !d.Write {
		c.Ui.Info(fmt.Sprintf("=== dry run: proposed changes to %s (pass -write to apply) ===", path))
		c.printDiff(path, newSrc)
		return nil
	}
	if err := os.WriteFile(path, newSrc, 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", path, err)
	}
	if err := helpers.GoImports(path); err != nil {
		return fmt.Errorf("formatting %q: %w", path, err)
	}
	c.Ui.Info(fmt.Sprintf("upgraded %s", path))
	return nil
}

// generateList renders and writes (or previews) the list resource and its test.
func (c UpgradeCommand) generateList(d *upgradeData, res *ir.ResourceIR) error {
	if !res.IsListable {
		return fmt.Errorf("resource %q has no subscription/resource-group or parent list operations", res.TerraformType)
	}
	listGo, err := gen.GenerateList(res)
	if err != nil {
		return err
	}
	listTestGo, err := gen.GenerateListTest(res)
	if err != nil {
		return err
	}

	dir := filepath.Dir(d.File)
	base := strings.TrimSuffix(filepath.Base(d.File), "_resource.go")
	listFile := filepath.Join(dir, base+"_resource_list.go")
	listTestFile := filepath.Join(dir, base+"_resource_list_test.go")

	if err := c.writeOrPreview(d, listFile, listGo); err != nil {
		return err
	}
	if err := c.writeOrPreview(d, listTestFile, listTestGo); err != nil {
		return err
	}

	regPath := filepath.Join(dir, "registration.go")
	return c.registerListResource(d, regPath, res.Name+"ListResource")
}

// registerListResource adds the generated list resource to the service package's
// registration.go ListResources() method, previewing or writing per -write.
func (c UpgradeCommand) registerListResource(d *upgradeData, regPath, listStruct string) error {
	if _, err := os.Stat(regPath); err != nil {
		c.Ui.Warn(fmt.Sprintf("no registration.go found at %s; skipping list resource registration", regPath))
		return nil
	}
	newSrc, changed, err := list_upgrade.RegisterListResource(regPath, listStruct)
	if err != nil {
		return fmt.Errorf("registering list resource: %w", err)
	}
	if !changed {
		c.Ui.Info(fmt.Sprintf("%s: %s already registered", regPath, listStruct))
		return nil
	}
	if !d.Write {
		c.Ui.Info(fmt.Sprintf("=== dry run: proposed changes to %s (pass -write to apply) ===", regPath))
		c.printDiff(regPath, newSrc)
		return nil
	}
	if err := os.WriteFile(regPath, newSrc, 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", regPath, err)
	}
	if err := helpers.GoImports(regPath); err != nil {
		return fmt.Errorf("formatting %q: %w", regPath, err)
	}
	c.Ui.Info(fmt.Sprintf("registered %s in %s", listStruct, regPath))
	return nil
}

func (c UpgradeCommand) writeOrPreview(d *upgradeData, path, content string) error {
	if !d.Write {
		c.Ui.Info(fmt.Sprintf("=== dry run: would generate %s ===", path))
		if formatted, err := format.Source([]byte(content)); err == nil {
			c.Ui.Output(string(formatted))
		} else {
			c.Ui.Output(content)
		}
		return nil
	}
	if !d.Overwrite {
		if _, err := os.Stat(path); err == nil {
			c.Ui.Warn(fmt.Sprintf("skipping existing file %s (use -overwrite to replace)", path))
			return nil
		}
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", path, err)
	}
	if err := helpers.GoImports(path); err != nil {
		return fmt.Errorf("formatting %q: %w", path, err)
	}
	c.Ui.Info(fmt.Sprintf("generated %s", path))
	return nil
}

// resolveIR resolves the resource IR from the Pandora Data API, using the
// explicit service/resource when supplied, otherwise the ARM type.
func (c UpgradeCommand) resolveIR(d *upgradeData, res *list_upgrade.Resource) (*ir.ResourceIR, error) {
	pandoraURL := d.PandoraURL
	if pandoraURL == "" && config != nil {
		pandoraURL = config.SchemaAPIURL
	}
	provider, org := d.provider(), d.org()
	client := pandora.NewClient(pandoraURL)

	name := d.Name
	if name == "" {
		name = strings.TrimPrefix(res.TerraformType, provider+"_")
	}
	return ir.Resolve(client, ir.Options{
		ARMType:           d.ARMType,
		Service:           d.Service,
		Resource:          d.Resource,
		APIVersion:        d.APIVersion,
		Name:              name,
		GoName:            d.GoName,
		ServicePackage:    res.Package,
		ProviderName:      provider,
		ProviderGithubOrg: org,
	})
}

func (d *upgradeData) provider() string {
	if d.Provider != "" {
		return d.Provider
	}
	if config != nil && config.ProviderName != "" && config.ProviderName != "commands" {
		return config.ProviderName
	}
	return "azurerm"
}

func (d *upgradeData) org() string {
	if d.Org != "" {
		return d.Org
	}
	if config != nil && config.ProviderGithubOrg != "" {
		return config.ProviderGithubOrg
	}
	return "hashicorp"
}

// reconcileIR overrides the Pandora-derived naming, package, model, ID and
// client fields with the values actually used by the hand-written resource, so
// the generated list code references real symbols.
func reconcileIR(res *ir.ResourceIR, a *list_upgrade.Resource) {
	res.Name = strings.TrimSuffix(a.StructName, "Resource")
	if a.SDKPackage != "" {
		res.SDKPackage = a.SDKPackage
	}
	if a.SDKImportPath != "" {
		res.SDKImportPath = a.SDKImportPath
	}
	if a.ClientField != "" {
		res.ClientField = a.ClientField
	}
	if a.ServiceName != "" {
		res.ServiceName = a.ServiceName
	}
	if a.Package != "" {
		res.ServicePackage = a.Package
	}
	if a.IDParseFunc != "" {
		res.IDParseFunc = a.IDParseFunc
	}
	if a.IDTypeName != "" {
		res.IDTypeName = a.IDTypeName
	}
	if a.TerraformType != "" {
		res.TerraformType = a.TerraformType
	}
}

func (c UpgradeCommand) printPlan(path string, p list_upgrade.Plan) {
	c.Ui.Info(fmt.Sprintf("Analyzing %s", path))
	c.Ui.Info(fmt.Sprintf("  kind:        %s", p.Kind))
	c.Ui.Info(fmt.Sprintf("  identity:    %s", present(p.HasIdentity)))
	c.Ui.Info(fmt.Sprintf("  flatten:     %s", present(p.HasFlatten)))
	if !p.Supported {
		c.Ui.Warn(fmt.Sprintf("  unsupported: %s", p.Reason))
	}
}

func present(v bool) string {
	if v {
		return "present"
	}
	return "missing"
}

// printDiff prints a unified diff between the on-disk file and the proposed
// (gofmt-formatted) content. Imports are resolved only on -write, so the preview
// reflects structural changes rather than final import ordering.
func (c UpgradeCommand) printDiff(path string, proposed []byte) {
	formatted := proposed
	if out, err := format.Source(proposed); err == nil {
		formatted = out
	}
	tmp, err := os.CreateTemp("", "scaff-upgrade-*.go")
	if err != nil {
		c.Ui.Output(string(formatted))
		return
	}
	defer os.Remove(tmp.Name())
	_, _ = tmp.Write(formatted)
	_ = tmp.Close()

	out, _ := exec.Command("diff", "-u", path, tmp.Name()).CombinedOutput()
	if len(out) == 0 {
		c.Ui.Output("(no textual differences)")
		return
	}
	c.Ui.Output(string(out))
}

func (c UpgradeCommand) Help() string {
	return `
Usage: scaff upgrade -file <path> [options]
       scaff upgrade -input <file.hcl>

Upgrades an existing typed resource so it can support List. It adds Resource
Identity and refactors Read into a reusable flatten method when those are
missing, then optionally generates the list resource.

Files are not modified without -write; by default the command performs a dry run
and prints the proposed changes.

Examples:
  # Report what an upgrade would do (read only):
  $ scaff upgrade -file internal/services/monitor/monitor_workspace_resource.go

  # Make a resource list-ready and generate the list resource:
  $ scaff upgrade -file internal/services/monitor/monitor_workspace_resource.go \
      -service Monitor -resource AzureMonitorWorkspaces -list -write

  # Upgrade one or more resources described in an HCL file:
  $ scaff upgrade -input internal/tools/scaff/examples/upgrade.hcl
`
}

func (c UpgradeCommand) Synopsis() string {
	return "Upgrades an existing resource to support List (adds identity and flatten)."
}
