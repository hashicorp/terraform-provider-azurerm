// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package commands

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	typed_upgrade "github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/typed-upgrade"
	"github.com/mitchellh/cli"
)

// TypedUpgradeCommand converts a native Plugin SDK resource
// (func resourceX() *pluginsdk.Resource) to the Typed SDK wrapper
// (type XResource struct{} implementing sdk.Resource).
type TypedUpgradeCommand struct {
	Ui cli.Ui
}

var _ cli.Command = TypedUpgradeCommand{}

type typedUpgradeData struct {
	File               string
	TerraformType      string
	ARMType            string
	Service            string
	Resource           string
	APIVersion         string
	PandoraURL         string
	Write              bool
	Overwrite          bool
	UpdateRegistration bool
	RegistrationPath   string
}

func (d *typedUpgradeData) parseArgs(args []string) error {
	set := flag.NewFlagSet("typed-upgrade", flag.ExitOnError)
	set.StringVar(&d.File, "file", "", "(Required) path to the existing untyped resource .go file to convert")
	set.StringVar(&d.TerraformType, "terraform-type", "", "(Optional) Terraform resource type override, e.g. \"azurerm_maps_creator\"; derived from the constructor name when omitted")
	set.StringVar(&d.ARMType, "arm-type", "", "(Optional) ARM resource type, e.g. \"Microsoft.Maps/accounts\"; enables Pandora IR resolution for nested block expand/flatten generation")
	set.StringVar(&d.Service, "service", "", "(Optional) explicit Pandora service name; overrides the value derived from -arm-type")
	set.StringVar(&d.Resource, "resource", "", "(Optional) explicit Pandora resource key; overrides the value derived from -arm-type")
	set.StringVar(&d.APIVersion, "api-version", "", "(Optional) API version; defaults to the latest non-preview version")
	set.StringVar(&d.PandoraURL, "pandora-url", "", "(Optional) Pandora Data API base URL; defaults to http://localhost:8080")
	set.BoolVar(&d.Write, "write", false, "(Optional) write the generated typed resource to disk; without it the command performs a dry run")
	set.BoolVar(&d.Overwrite, "overwrite", false, "(Optional) overwrite an existing generated file")
	set.BoolVar(&d.UpdateRegistration, "update-registration", true, "(Optional) add the typed resource to Resources() and remove it from SupportedResources() in registration.go")
	set.StringVar(&d.RegistrationPath, "registration", "", "(Optional) explicit path to registration.go; derived from the resource file when omitted")

	if err := set.Parse(args); err != nil {
		return err
	}
	if d.File == "" {
		return errors.New("-file is required")
	}
	return nil
}

func (c TypedUpgradeCommand) Help() string {
	return `Usage: scaff typed-upgrade [options]

  Converts a native Plugin SDK resource (func resourceX() *pluginsdk.Resource)
  to the Typed SDK wrapper (type XResource struct{} implementing sdk.Resource).

  The command performs a dry run by default, printing the plan and a preview of
  the generated file. Pass -write to apply the changes.

  Pass -arm-type with a running Pandora Data API to generate typed expand/flatten
  helpers for nested block fields (see 'serve-data-api.sh').

Options:

  -file <path>              (Required) path to the existing untyped resource file.
  -terraform-type <string>  Override the derived Terraform resource type string.
  -arm-type <string>        ARM resource type, e.g. "Microsoft.Maps/accounts"; enables
                            nested block expand/flatten generation via Pandora.
  -service <string>         Explicit Pandora service name (overrides -arm-type derivation).
  -resource <string>        Explicit Pandora resource key.
  -api-version <string>     API version; defaults to the latest non-preview version.
  -pandora-url <string>     Pandora Data API URL (default: http://localhost:8080).
  -write                    Write changes to disk (default: dry run).
  -overwrite                Overwrite existing files.
  -update-registration      Update registration.go (default: true).
  -registration <path>      Explicit path to registration.go.
`
}

func (c TypedUpgradeCommand) Synopsis() string {
	return "Convert a native Plugin SDK resource to the Typed SDK wrapper"
}

func (c TypedUpgradeCommand) Run(args []string) int {
	d := &typedUpgradeData{}
	if err := d.parseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Analysing %s ...", d.File))
	info, err := typed_upgrade.Analyze(d.File)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("analysis failed: %+v", err))
		return 1
	}

	plan := typed_upgrade.PlanUpgrade(info)
	c.printPlan(d.File, info, plan)

	if !plan.Supported {
		c.Ui.Error(fmt.Sprintf("cannot upgrade: %s", plan.Reason))
		return 1
	}

	if !d.Write {
		c.Ui.Info("\nDry run complete. Pass -write to apply the changes.")
		return 0
	}

	opts := typed_upgrade.UpgradeOptions{
		Write:              d.Write,
		Overwrite:          d.Overwrite,
		UpdateRegistration: d.UpdateRegistration,
		RegistrationPath:   d.RegistrationPath,
		TerraformType:      d.TerraformType,
		PandoraURL:         d.PandoraURL,
		ARMType:            d.ARMType,
		Service:            d.Service,
		Resource:           d.Resource,
		APIVersion:         d.APIVersion,
	}

	res, err := typed_upgrade.Upgrade(info, opts)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("upgrade failed: %+v", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("✓ Generated typed resource: %s", res.GeneratedPath))
	c.Ui.Info(fmt.Sprintf("✓ Original renamed to:      %s", res.RenamedPath))

	// Run goimports on the generated file.
	if path, err := exec.LookPath("goimports"); err == nil {
		cmd := exec.Command(path, "-w", res.GeneratedPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			c.Ui.Warn(fmt.Sprintf("goimports: %s\n%s", err, out))
		} else {
			c.Ui.Info(fmt.Sprintf("✓ goimports applied to %s", res.GeneratedPath))
		}
	}

	c.Ui.Info("\nNext steps:")
	c.Ui.Info("  1. Review the generated file and replace metadata.ResourceData.Get/Set calls with model fields.")
	c.Ui.Info("  2. Review the _legacy.go file and delete it once all logic has been moved.")
	c.Ui.Info("  3. Run: go vet ./internal/services/<service>/...")
	c.Ui.Info("  4. Run: go test ./internal/services/<service>/...")

	return 0
}

func (c TypedUpgradeCommand) printPlan(file string, info *typed_upgrade.Info, plan typed_upgrade.Plan) {
	c.Ui.Output(fmt.Sprintf("\n=== Typed Upgrade Plan: %s ===\n", file))
	c.Ui.Output(fmt.Sprintf("  Terraform type : %s", plan.TerraformType))
	c.Ui.Output(fmt.Sprintf("  Struct name    : %s", plan.StructName))
	c.Ui.Output(fmt.Sprintf("  Model name     : %s", plan.ModelName))
	c.Ui.Output(fmt.Sprintf("  Has Update     : %v", plan.HasUpdate))
	c.Ui.Output(fmt.Sprintf("  Schema fields  : %d", len(info.Fields)))
	c.Ui.Output(fmt.Sprintf("  SDK package    : %s", info.SDKPackage))
	c.Ui.Output(fmt.Sprintf("  ID type        : %s", info.IDTypeName))
	if info.PandoraIR != nil {
		c.Ui.Output(fmt.Sprintf("  Pandora IR     : resolved (api %s, %d blocks matched)",
			info.PandoraIR.APIVersion, typed_upgrade.IRBlockCount(info)))
	}
	if len(plan.Warnings) > 0 {
		c.Ui.Output("\nWarnings:")
		for _, w := range plan.Warnings {
			c.Ui.Warn(fmt.Sprintf("  ! %s", w))
		}
	}
	c.Ui.Output("")
	if !plan.Supported {
		return
	}
	actions := []string{
		fmt.Sprintf("generate %s (typed resource)", plan.StructName),
		fmt.Sprintf("rename original to *_legacy.go"),
	}
	c.Ui.Output("Actions:")
	for _, a := range actions {
		c.Ui.Output(fmt.Sprintf("  + %s", a))
	}
	c.Ui.Output(fmt.Sprintf("\nField breakdown: %s", fieldBreakdown(info)))
}

func fieldBreakdown(info *typed_upgrade.Info) string {
	var args, attrs []string
	for _, f := range info.Fields {
		if f.IsAttribute() {
			attrs = append(attrs, f.TFName)
		} else {
			args = append(args, f.TFName)
		}
	}
	parts := []string{}
	if len(args) > 0 {
		parts = append(parts, fmt.Sprintf("arguments=[%s]", strings.Join(args, ", ")))
	}
	if len(attrs) > 0 {
		parts = append(parts, fmt.Sprintf("attributes=[%s]", strings.Join(attrs, ", ")))
	}
	if len(parts) == 0 {
		return "(none)"
	}
	return strings.Join(parts, " ")
}
