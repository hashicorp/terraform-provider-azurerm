// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
	list_upgrade "github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/list-upgrade"
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
	ListMethod    string
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
	ListMethod    string `hcl:"list_method,optional"`
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
	set.StringVar(&d.ListMethod, "list-method", "", "(Optional) parent-scoped SDK list method name (without the Complete suffix); overrides the value derived from the Get method")
	set.StringVar(&d.IdentityProps, "identity-properties", "", "(Optional) -properties value for the identity test go:generate directive; defaults to \"name,resource_group_name\"")
	set.BoolVar(&d.Identity, "identity", false, "(Optional) add Resource Identity if missing")
	set.BoolVar(&d.Flatten, "flatten", false, "(Optional) refactor Read into a reusable flatten method if missing")
	set.BoolVar(&d.List, "list", true, "(Optional) make the resource list-ready (identity + flatten) and generate the list resource")
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

	var failures []string
	for _, b := range file.Resources {
		if b.File == "" {
			c.Ui.Error(fmt.Sprintf("resource %q: file is required", b.Name))
			failures = append(failures, b.Name)
			continue
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
			ListMethod:    b.ListMethod,
			IdentityProps: b.IdentityProps,
			Identity:      derefBool(b.Identity),
			Flatten:       derefBool(b.Flatten),
			// List defaults to true (matching the -list CLI flag) so a block need
			// only set `list = false` to opt out.
			List:       b.List == nil || *b.List,
			Write:      write,
			Overwrite:  overwrite,
			PandoraURL: file.PandoraURL,
			Provider:   file.Provider,
			Org:        file.Org,
		}
		// Continue past a failed resource so one unsupported entry (e.g. an
		// untyped resource) does not abort the whole batch; failures are
		// collected and reported together at the end.
		if err := c.run(rd); err != nil {
			c.Ui.Error(fmt.Sprintf("resource %q: %+v", b.Name, err))
			failures = append(failures, b.Name)
			continue
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%d of %d resource(s) could not be upgraded: %s", len(failures), len(file.Resources), strings.Join(failures, ", "))
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

	// A resource whose list operations are derivable from source — a detected
	// parent scope, or subscription/resource-group list methods found in the
	// vendored SDK — needs no Pandora resolution.
	sourceOnly := res.ParentIDType != "" || res.ListSubscriptionMethod != "" || res.ListResourceGroupMethod != ""

	// Read model: -read-model flag > vendor-derived (authoritative: the type of
	// resp.Model) > Pandora (last resort). The flatten parameter and list results
	// must match resp.Model exactly.
	readModel := d.ReadModel
	if readModel == "" {
		readModel = res.ReadModel
	}

	// Resolve the IR from Pandora only when we still need the read model or the
	// list operations (i.e. not fully source-derivable).
	var resolved *ir.ResourceIR
	if !sourceOnly && (d.List || (wantFlatten && readModel == "")) {
		resolved, err = c.resolveIR(d, res)
		if err != nil {
			return fmt.Errorf("resolving API definitions (provide -read-model to skip, or -arm-type/-service+-resource): %w", err)
		}
	}

	if readModel == "" && resolved != nil {
		readModel = resolved.ReadModel
	}
	if wantFlatten && readModel == "" {
		return errors.New("refactoring Read into flatten needs the SDK read model; provide -read-model (required for parent-scoped resources) or -arm-type")
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
		if sourceOnly {
			if readModel == "" {
				return errors.New("-list requires the SDK read model; it could not be derived from the vendored SDK — provide -read-model")
			}
			resolved = c.sourceListIR(res, d, readModel)
		} else {
			if resolved == nil {
				return errors.New("-list requires API definitions; provide -arm-type or -service and -resource")
			}
			reconcileIR(resolved, res)
			// The vendored-SDK read model (the type of resp.Model) is
			// authoritative; override Pandora's classification so the list
			// results slice (var results []<pkg>.<model>) matches the flatten.
			if readModel != "" {
				resolved.ReadModel = readModel
			}
		}
		if err := c.generateList(d, resolved); err != nil {
			return err
		}
	}
	return nil
}

// sourceListIR builds the list IR for an untyped, parent-scoped child resource
// entirely from the parsed source plus the supplied read model, bypassing
// Pandora (whose grouping cannot always distinguish such sub-resources).
func (c UpgradeCommand) sourceListIR(a *list_upgrade.Resource, d *upgradeData, readModel string) *ir.ResourceIR {
	provider, org := d.provider(), d.org()
	name := d.Name
	if name == "" {
		name = a.BaseName
	}
	listMethod := a.ListMethod
	if d.ListMethod != "" {
		listMethod = d.ListMethod
	}
	terraformType := a.TerraformType
	if terraformType == "" {
		terraformType = provider + "_" + name
	}
	res := &ir.ResourceIR{
		ProviderName:      provider,
		ProviderGithubOrg: org,
		ServicePackage:    a.Package,
		ServiceName:       a.ServiceName,
		ClientField:       a.ClientField,
		TerraformType:     terraformType,
		SDKPackage:        a.SDKPackage,
		SDKImportPath:     a.SDKImportPath,
		ReadModel:         readModel,
		IDPackage:         a.IDPackage,
		IDImportPath:      a.IDImportPath,
		IDParseFunc:       a.IDParseFunc,
		IDTypeName:        a.IDTypeName,
		IsListable:        true,
	}
	if a.ParentIDType != "" {
		// Parent-scoped child: a single list call under the parent ID.
		res.ListByParentOp = listMethod
		res.ParentIDType = a.ParentIDType
		res.ParentListAttr = a.ParentAttr
		res.ParentParseFunc = a.ParentParseFunc
		res.ParentValidateFunc = a.ParentValidateFunc
		res.ParentPackage = a.ParentPackage
		res.ParentImportPath = a.ParentImportPath
	} else {
		// Top-level: subscription/resource-group list methods (from the SDK).
		res.ListBySubscriptionOp = a.ListSubscriptionMethod
		res.ListByResourceGroupOp = a.ListResourceGroupMethod
	}
	reconcileParentListOp(res, a.Path)
	setListOptionFlags(res, a.Path)
	if a.Kind == list_upgrade.KindUntyped {
		res.Untyped = true
		res.Name = a.BaseName
		res.ConstructorFunc = a.ConstructorFunc
		res.FlattenFunc = a.ConstructorFunc + "Flatten"
		res.FlattenIDValue = a.FlattenIDValue
		res.FlattenNeedsContext = a.FlattenNeedsContext
		res.FlattenClientType = a.ClientTypeName
	} else {
		res.Name = strings.TrimSuffix(a.StructName, "Resource")
	}
	return res
}

// setListOptionFlags resolves, from the vendored SDK, whether each list Complete
// method takes a trailing options argument so the generator emits it. anchorPath
// is any path inside the provider module (used to locate vendor/).
func setListOptionFlags(res *ir.ResourceIR, anchorPath string) {
	res.ListBySubscriptionHasOptions = list_upgrade.ListMethodTakesOptions(anchorPath, res.SDKImportPath, res.ListBySubscriptionOp)
	res.ListByResourceGroupHasOptions = list_upgrade.ListMethodTakesOptions(anchorPath, res.SDKImportPath, res.ListByResourceGroupOp)
	res.ListByParentHasOptions = list_upgrade.ListMethodTakesOptions(anchorPath, res.SDKImportPath, res.ListByParentOp)
}

// reconcileParentListOp reconciles a parent-scoped list operation name against
// the vendored SDK. Pandora (and the typed source-derived path) may name it
// differently from the SDK (e.g. "List" vs "ListByMongoCluster"); when the
// resolved name has no matching SDK method, it is replaced by the SDK method
// whose id parameter is the parent id type.
func reconcileParentListOp(res *ir.ResourceIR, anchorPath string) {
	if res.ListByParentOp == "" || res.ParentIDType == "" {
		return
	}
	if list_upgrade.ListCompleteMethodExists(anchorPath, res.SDKImportPath, res.ListByParentOp) {
		return
	}
	if m := list_upgrade.ParentListMethodByID(anchorPath, res.SDKImportPath, res.ParentIDType); m != "" {
		res.ListByParentOp = m
	}
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

	dir := filepath.Dir(d.File)
	base := strings.TrimSuffix(filepath.Base(d.File), "_resource.go")
	listFile := filepath.Join(dir, base+"_resource_list.go")

	if err := c.writeOrPreview(d, listFile, listGo); err != nil {
		return err
	}

	if d.Write {
		generateListDoc(c.Ui, listFile, d.Overwrite)
	}

	// Generate the acceptance test, referencing the resource's existing test
	// struct (whose casing may differ from the resource base name).
	if res.TestStructName == "" {
		res.TestStructName = list_upgrade.TestStructName(d.File)
	}
	if res.TestStructName != "" {
		listTestGo, err := gen.GenerateListTest(res)
		if err != nil {
			return err
		}
		listTestFile := filepath.Join(dir, base+"_resource_list_test.go")
		if err := c.writeOrPreview(d, listTestFile, listTestGo); err != nil {
			return err
		}
	} else {
		c.Ui.Warn(fmt.Sprintf("no acceptance-test struct found alongside %s; skipping list test generation", d.File))
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
	if a.Kind == list_upgrade.KindUntyped {
		res.Untyped = true
		res.Name = a.BaseName
		res.ConstructorFunc = a.ConstructorFunc
		res.FlattenFunc = a.ConstructorFunc + "Flatten"
		res.FlattenIDValue = a.FlattenIDValue
		res.FlattenNeedsContext = a.FlattenNeedsContext
		res.FlattenClientType = a.ClientTypeName
		res.IDPackage = a.IDPackage
		res.IDImportPath = a.IDImportPath
	} else {
		res.Name = strings.TrimSuffix(a.StructName, "Resource")
	}
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
	reconcileParentListOp(res, a.Path)
	setListOptionFlags(res, a.Path)
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
