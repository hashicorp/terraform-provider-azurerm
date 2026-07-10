package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/gen"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"

	"github.com/mitchellh/cli"
)

type GenerateCommand struct {
	Ui cli.Ui
}

var _ cli.Command = GenerateCommand{}

type GenerateData struct {
	ARMType        string
	Service        string
	Resource       string
	APIVersion     string
	Name           string
	GoName         string
	ServicePackage string
	PandoraURL     string
	OutputPath     string
	Provider       string
	Org            string
	InputFile      string
	Overwrite      bool
	List           bool
	GenResource    bool
	DataSource     bool
	Config         *helpers.Configuration
}

// generateFile is the schema of an HCL file supplied via -input. It describes
// one or more resources to generate, plus optional global overrides that apply
// to every resource in the file.
type generateFile struct {
	PandoraURL string          `hcl:"pandora_url,optional"`
	OutputPath string          `hcl:"output_path,optional"`
	Provider   string          `hcl:"provider,optional"`
	Org        string          `hcl:"org,optional"`
	Overwrite  *bool           `hcl:"overwrite,optional"`
	Resources  []resourceBlock `hcl:"resource,block"`
}

// resourceBlock is a single "resource" block within a generate input file. The
// block label is the terraform resource name (snake_case, without the provider
// prefix), e.g. resource "redhat_openshift_cluster" { ... }.
type resourceBlock struct {
	Name           string `hcl:"name,label"`
	ARMType        string `hcl:"arm_type,optional"`
	Service        string `hcl:"service,optional"`
	Resource       string `hcl:"pandora_resource,optional"`
	APIVersion     string `hcl:"api_version,optional"`
	GoName         string `hcl:"go_name,optional"`
	ServicePackage string `hcl:"servicepackage,optional"`
	OutputPath     string `hcl:"path,optional"`
	GenResource    *bool  `hcl:"gen_resource,optional"`
	List           bool   `hcl:"list,optional"`
	DataSource     bool   `hcl:"data_source,optional"`
}

// resourceSpec is a single, fully-resolved generation request shared by the
// flag- and file-driven code paths.
type resourceSpec struct {
	ARMType        string
	Service        string
	Resource       string
	APIVersion     string
	Name           string
	GoName         string
	ServicePackage string
	OutputPath     string
	GenResource    bool
	List           bool
	DataSource     bool
}

func (s resourceSpec) validate() error {
	if s.Name == "" {
		return errors.New("a resource name (block label) is required")
	}
	if s.ARMType == "" && (s.Service == "" || s.Resource == "") {
		return errors.New("arm_type is required (or provide both service and pandora_resource)")
	}
	if !s.GenResource && !s.List && !s.DataSource {
		return errors.New("nothing to generate: enable gen_resource (default), list and/or data_source")
	}
	return nil
}

func (d *GenerateData) parseArgs(args []string) (errs []error) {
	set := flag.NewFlagSet("generate", flag.ExitOnError)
	set.StringVar(&d.ARMType, "arm-type", "", "(Required) ARM resource type, e.g. \"Microsoft.RedHatOpenShift/openShiftClusters\"")
	set.StringVar(&d.Name, "name", "", "(Required) terraform resource name (snake_case, without provider prefix), e.g. \"redhat_openshift_cluster\"")
	set.StringVar(&d.GoName, "go-name", "", "(Optional) Go identifier for the resource, e.g. \"RedHatOpenShiftCluster\"; derived from -name when omitted")
	set.StringVar(&d.ServicePackage, "servicepackage", "", "(Optional) service package directory; derived from the service name when omitted")
	set.StringVar(&d.APIVersion, "api-version", "", "(Optional) API version; defaults to the latest non-preview version")
	set.StringVar(&d.Service, "service", "", "(Optional) explicit Pandora service name; overrides the value derived from -arm-type")
	set.StringVar(&d.Resource, "resource", "", "(Optional) explicit Pandora resource key; overrides the value derived from -arm-type")
	set.StringVar(&d.PandoraURL, "pandora-url", "", "(Optional) Pandora Data API base URL; defaults to the configured schema_api_url")
	set.StringVar(&d.OutputPath, "path", "", "(Optional) output directory; defaults to {service_packages_path}/{servicepackage}")
	set.StringVar(&d.Provider, "provider", "", "(Optional) provider name override, e.g. \"azurerm\"")
	set.StringVar(&d.Org, "org", "", "(Optional) provider GitHub org override, e.g. \"hashicorp\"")
	set.BoolVar(&d.List, "list", false, "(Optional) also generate a list resource and its acceptance test")
	set.BoolVar(&d.GenResource, "gen-resource", true, "(Optional) generate the resource file; set to false with -list to generate only the list resource (assumes the resource already exists with a flatten method)")
	set.BoolVar(&d.DataSource, "data-source", false, "(Optional) also generate a data source (shares nested block structs with the resource)")
	set.StringVar(&d.InputFile, "input", "", "(Optional) path to an HCL file describing one or more resources to generate; when set, the per-resource flags are ignored")
	set.BoolVar(&d.Overwrite, "overwrite", false, "(Optional) overwrite existing generated files; by default existing files are left untouched")

	if err := set.Parse(args); err != nil {
		return []error{err}
	}

	// In file mode the resource inputs are read from the HCL file, so the
	// per-resource flags are neither required nor consulted.
	if d.InputFile != "" {
		return nil
	}

	if d.ARMType == "" && (d.Service == "" || d.Resource == "") {
		errs = append(errs, errors.New("-arm-type is required (or provide both -service and -resource)"))
	}
	if d.Name == "" {
		errs = append(errs, errors.New("-name is required"))
	}
	if !d.GenResource && !d.List && !d.DataSource {
		errs = append(errs, errors.New("nothing to generate: enable -gen-resource (default), -list and/or -data-source"))
	}
	return errs
}

func (c GenerateCommand) Run(args []string) int {
	data := &GenerateData{Config: config}
	if err := data.parseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generate(c.Ui); err != nil {
		c.Ui.Error(fmt.Sprintf("generating resource: %+v", err))
		return 1
	}
	return 0
}

// generate dispatches to file- or flag-driven generation depending on whether
// an -input HCL file was supplied.
func (d GenerateData) generate(ui cli.Ui) error {
	if d.InputFile != "" {
		return d.generateFromFile(ui)
	}
	return d.generateFromFlags(ui)
}

// generateFromFlags generates a single resource from the command-line flags,
// falling back to .scaff.hcl config defaults.
func (d GenerateData) generateFromFlags(ui cli.Ui) error {
	pandoraURL := d.PandoraURL
	if pandoraURL == "" && d.Config != nil {
		pandoraURL = d.Config.SchemaAPIURL
	}

	provider, org := d.configProviderOrg()
	if d.Provider != "" {
		provider = d.Provider
	}
	if d.Org != "" {
		org = d.Org
	}

	client := pandora.NewClient(pandoraURL)
	spec := resourceSpec{
		ARMType:        d.ARMType,
		Service:        d.Service,
		Resource:       d.Resource,
		APIVersion:     d.APIVersion,
		Name:           d.Name,
		GoName:         d.GoName,
		ServicePackage: d.ServicePackage,
		OutputPath:     d.OutputPath,
		GenResource:    d.GenResource,
		List:           d.List,
		DataSource:     d.DataSource,
	}
	return d.generateResource(ui, client, spec, provider, org, d.baseOutputPath(""))
}

// generateFromFile reads an HCL input file and generates every resource block it
// declares. Global attributes in the file override the .scaff.hcl config
// defaults; the per-resource CLI flags are ignored in this mode.
func (d GenerateData) generateFromFile(ui cli.Ui) error {
	var file generateFile
	if err := hclsimple.DecodeFile(d.InputFile, nil, &file); err != nil {
		return fmt.Errorf("reading input file %q: %w", d.InputFile, err)
	}
	if len(file.Resources) == 0 {
		return fmt.Errorf("input file %q defines no resource blocks", d.InputFile)
	}

	pandoraURL := file.PandoraURL
	if pandoraURL == "" && d.Config != nil {
		pandoraURL = d.Config.SchemaAPIURL
	}

	provider, org := d.configProviderOrg()
	if file.Provider != "" {
		provider = file.Provider
	}
	if file.Org != "" {
		org = file.Org
	}
	if file.Overwrite != nil {
		d.Overwrite = *file.Overwrite
	}
	baseOutput := d.baseOutputPath(file.OutputPath)

	client := pandora.NewClient(pandoraURL)
	for _, rb := range file.Resources {
		spec := resourceSpec{
			ARMType:        rb.ARMType,
			Service:        rb.Service,
			Resource:       rb.Resource,
			APIVersion:     rb.APIVersion,
			Name:           rb.Name,
			GoName:         rb.GoName,
			ServicePackage: rb.ServicePackage,
			OutputPath:     rb.OutputPath,
			GenResource:    rb.GenResource == nil || *rb.GenResource,
			List:           rb.List,
			DataSource:     rb.DataSource,
		}
		if err := spec.validate(); err != nil {
			return fmt.Errorf("resource %q: %w", rb.Name, err)
		}
		if err := d.generateResource(ui, client, spec, provider, org, baseOutput); err != nil {
			return fmt.Errorf("resource %q: %w", rb.Name, err)
		}
	}
	return nil
}

// configProviderOrg returns the provider name and GitHub org from the loaded
// config, or empty strings when no config is present.
func (d GenerateData) configProviderOrg() (provider, org string) {
	if d.Config != nil {
		return d.Config.ProviderName, d.Config.ProviderGithubOrg
	}
	return "", ""
}

// baseOutputPath resolves the base output directory, preferring the supplied
// override, then the configured service_packages_path, then a sensible default.
func (d GenerateData) baseOutputPath(override string) string {
	if override != "" {
		return override
	}
	if d.Config != nil && d.Config.ServicePackagesPath != "" {
		return d.Config.ServicePackagesPath
	}
	return "internal/services"
}

// generateResource resolves a single resource from Pandora and writes the
// requested resource, list and/or data source files.
func (d GenerateData) generateResource(ui cli.Ui, client *pandora.Client, spec resourceSpec, provider, org, baseOutput string) error {
	res, err := ir.Resolve(client, ir.Options{
		ARMType:           spec.ARMType,
		Service:           spec.Service,
		Resource:          spec.Resource,
		APIVersion:        spec.APIVersion,
		Name:              spec.Name,
		GoName:            spec.GoName,
		ServicePackage:    spec.ServicePackage,
		ProviderName:      provider,
		ProviderGithubOrg: org,
	})
	if err != nil {
		return err
	}

	outputDir := spec.OutputPath
	if outputDir == "" {
		outputDir = filepath.Join(baseOutput, res.ServicePackage)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory %q: %w", outputDir, err)
	}

	if spec.GenResource {
		resourceGo, err := gen.Generate(res)
		if err != nil {
			return err
		}
		resourceFile := filepath.Join(outputDir, fmt.Sprintf("%s_resource.go", spec.Name))
		if err := d.writeGenerated(ui, resourceFile, resourceGo); err != nil {
			return err
		}
	}

	if spec.List {
		listGo, err := gen.GenerateList(res)
		if err != nil {
			return err
		}
		listFile := filepath.Join(outputDir, fmt.Sprintf("%s_resource_list.go", spec.Name))
		if err := d.writeGenerated(ui, listFile, listGo); err != nil {
			return err
		}

		listTestGo, err := gen.GenerateListTest(res)
		if err != nil {
			return err
		}
		listTestFile := filepath.Join(outputDir, fmt.Sprintf("%s_resource_list_test.go", spec.Name))
		if err := d.writeGenerated(ui, listTestFile, listTestGo); err != nil {
			return err
		}
	}

	if spec.DataSource {
		dsGo, err := gen.GenerateDataSource(res)
		if err != nil {
			return err
		}
		dsFile := filepath.Join(outputDir, fmt.Sprintf("%s_data_source.go", spec.Name))
		if err := d.writeGenerated(ui, dsFile, dsGo); err != nil {
			return err
		}
	}

	return nil
}

// writeGenerated writes and formats a generated file, reporting the result.
// Existing files are left untouched unless -overwrite was supplied, so the tool
// never silently clobbers hand-written resources.
func (d GenerateData) writeGenerated(ui cli.Ui, path, content string) error {
	if !d.Overwrite {
		if _, err := os.Stat(path); err == nil {
			ui.Warn(fmt.Sprintf("skipping existing file %s (use -overwrite to replace)", path))
			return nil
		}
	}
	if err := writeAndFormat(path, content); err != nil {
		return err
	}
	ui.Info(fmt.Sprintf("generated %s", path))
	return nil
}

func writeAndFormat(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", path, err)
	}
	if err := helpers.GoImports(path); err != nil {
		return fmt.Errorf("formatting %q: %w", path, err)
	}
	return nil
}

func (c GenerateCommand) Help() string {
	return `
Usage: scaff generate [options]

Generates a typed (internal/sdk) resource - schema, CRUD, and expand/flatten
functions - by reading the API definitions from a running Pandora Data API.

Existing files are left untouched by default; pass -overwrite to replace them.

Example:
$ scaff generate -arm-type="Microsoft.RedHatOpenShift/openShiftClusters" -name="redhat_openshift_cluster" -go-name="RedHatOpenShiftCluster" -servicepackage="redhatopenshift"

Example (list resource only - assumes the resource already exists with a flatten method):
$ scaff generate -arm-type="Microsoft.RedHatOpenShift/openShiftClusters" -name="redhat_openshift_cluster" -go-name="RedHatOpenShiftCluster" -servicepackage="redhatopenshift" -gen-resource=false -list

Example (resource + data source):
$ scaff generate -arm-type="Microsoft.RedHatOpenShift/openShiftClusters" -name="redhat_openshift_cluster" -go-name="RedHatOpenShiftCluster" -servicepackage="redhatopenshift" -data-source

Example (read one or more resources from an HCL file):
$ scaff generate -input="resources.hcl"

Where resources.hcl looks like:

  # optional global overrides (default to .scaff.hcl / built-in defaults)
  pandora_url = "http://localhost:8080"
  output_path = "internal/services"
  provider    = "azurerm"
  org         = "hashicorp"
  overwrite   = false                   # set true to replace existing files

  resource "redhat_openshift_cluster" {
    arm_type       = "Microsoft.RedHatOpenShift/openShiftClusters"
    go_name        = "RedHatOpenShiftCluster"
    servicepackage = "redhatopenshift"
    api_version    = "2025-07-25" # optional; defaults to latest non-preview
    list           = true         # optional; also generate a list resource
    data_source    = true         # optional; also generate a data source
    # gen_resource = false        # optional; defaults to true
    # path         = "..."        # optional; full output dir override
    # service          = "..."    # optional; overrides the service derived from arm_type
    # pandora_resource = "..."    # optional; overrides the resource key derived from arm_type
  }
`
}

func (c GenerateCommand) Synopsis() string {
	return "generates a typed resource from the Pandora Data API."
}
