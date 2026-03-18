package commands

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/templatehelpers"
	"github.com/mitchellh/cli"
)

type ListResourceCommand struct {
	Ui cli.Ui
}

var _ cli.Command = &ListResourceCommand{}

type ResourceInput struct {
	Service              string `json:"service"`
	Resource             string `json:"resource"`
	IncludeServiceInName bool   `json:"include_service_in_name"`
	FullParent           string `json:"full_parent"`
	Parent               string `json:"parent"`
	ParentTerraformName  string `json:"parent_terraform_name"`
	TerraformName        string `json:"terraform_name"`
	IDStructure          string `json:"id_structure"`
	Path                 string `json:"path"`
}

type ListResourceData struct {
	// Core
	ServiceName   string
	ResourceName  string
	PackageName   string
	ResourceLower string
	FullParent    string
	Parent        string
	ParentLower   string

	IdPackage string
	IdType    string
	IdName    string

	UseResourceGroup bool
	// Terraform
	TerraformResourceName       string
	ParentTerraformResourceName string

	// Go identifiers
	ResourceStruct     string
	ListResourceStruct string
	ModelStruct        string

	// Display
	ChildDisplayName string

	// Files
	OutputFile     string
	OutputTestFile string
	AzureSDKImport string
}

func loadResources(filePath string) ([]ResourceInput, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var resources []ResourceInput

	err = json.Unmarshal(data, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (c ListResourceCommand) Run(args []string) int {
	input := &ResourceInput{}

	if err := input.parseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := c.processResource(*input); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c ListResourceCommand) Synopsis() string {
	return "Generate list resource boilerplate code from JSON file or CLI arguments"
}

func (c ListResourceCommand) Help() string {
	return `
Usage: scaff list-resource [options]

  Generates list resource boilerplate code. Can accept either a JSON file
  with multiple resource definitions or individual CLI arguments for a single resource.

Options:
  -json=<path>                    Path to JSON file containing resource definitions (use instead of full cli)
  -service=<name>                 (Required) Service name (e.g., "PrivateDns")
  -resource=<name>                (Required) Resource name (e.g., "CNameRecord")
  -include_service_in_name=<bool> (Optional) Include service name in generated identifiers, defaults to false
  -full_parent=<name>             (Optional) Full parent resource name
  -parent=<name>                  (Optional) Parent resource name (e.g., "PrivateDnsZone"), defaults to resource_group
  -parent_terraform_name=<name>   (Optional) Parent Terraform resource name (e.g., "private_dns_zone")
  -terraform_name=<name>          (Optional) Terraform resource name (e.g., "private_dns_cname_record")
  -id_structure=<type>            (Required) ID structure type (e.g., "privatedns.RecordType")
  -path=<path>                    (Optional) Output path for generated files

Examples:
  # Using JSON file
  scaff list-resource -json=internal/tools/scaff/commands/input_example/listResources.json

  # Using CLI arguments
  scaff list-resource -service="PrivateDns" -resource="CNameRecord" -parent="PrivateDnsZone" -terraform_name="private_dns_cname_record" -id_structure="privatedns.RecordType" -path="internal/services/privatedns/"
`
}

func (input *ResourceInput) parseArgs(args []string) (errs []error) {
	argSet := flag.NewFlagSet("list-resource", flag.ContinueOnError)

	var jsonFile string
	argSet.StringVar(&jsonFile, "json", "", "Path to JSON file containing resource definitions")
	argSet.StringVar(&input.Service, "service", "", "Service name")
	argSet.StringVar(&input.Resource, "resource", "", "Resource name")
	argSet.BoolVar(&input.IncludeServiceInName, "include_service_in_name", false, "(Optional) Include service name in generated identifiers")
	argSet.StringVar(&input.FullParent, "full_parent", "", "(Optional) Full parent resource name, if different than parent")
	argSet.StringVar(&input.Parent, "parent", "", "(Optional) Parent resource name, will default to resource group if not specified")
	argSet.StringVar(&input.ParentTerraformName, "parent_terraform_name", "", "(Optional) Terraform resource name of the parent, will attempt to derive it if not provided")
	argSet.StringVar(&input.TerraformName, "terraform_name", "", "(Optional) Terraform resource name, will attempt to derive it if not provided")
	argSet.StringVar(&input.IDStructure, "id_structure", "", "ID structure type")
	argSet.StringVar(&input.Path, "path", "", "Output path for generated files")

	if err := argSet.Parse(args); err != nil {
		errs = append(errs, err)
		return
	}

	// If JSON file is provided, load from it
	if jsonFile != "" {
		resources, err := loadResources(jsonFile)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to load JSON file: %w", err))
			return
		}
		if len(resources) == 0 {
			errs = append(errs, errors.New("JSON file contains no resources"))
			return
		}
		*input = resources[0]
		return
	}

	// Validate required fields when using CLI arguments
	switch {
	case input.Service == "":
		errs = append(errs, errors.New("service is required"))
	case input.Resource == "":
		errs = append(errs, errors.New("resource is required"))
	case input.IDStructure == "":
		errs = append(errs, errors.New("id_structure is required"))
	case input.Path == "":
		input.Path = fmt.Sprintf(".")
	}

	return
}

func (c ListResourceCommand) processResource(input ResourceInput) error {
	data := derive(input)

	templateName := "list.go.gotpl"
	templatePath := "templates"
	outputPath := input.Path + data.OutputFile
	if err := data.renderTemplate(templatePath, templateName, outputPath); err != nil {
		return fmt.Errorf("failed to render %s: %w", outputPath, err)
	}
	c.Ui.Info(fmt.Sprintf("✅ generated %s", outputPath))

	templateName = "list_test.go.gotpl"
	testOutputPath := input.Path + data.OutputTestFile
	if err := data.renderTemplate(templatePath, templateName, testOutputPath); err != nil {
		return fmt.Errorf("failed to render %s: %w", testOutputPath, err)
	}
	c.Ui.Info(fmt.Sprintf("✅ generated %s", testOutputPath))

	return nil
}

func (data *ListResourceData) renderTemplate(templatePath, templateName, outputPath string) error {
	tmpl := template.Must(template.New(templateName).Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, fmt.Sprintf("%s/%s", templatePath, templateName)))

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, data)
}

func derive(input ResourceInput) *ListResourceData {
	service := input.Service   // Mssql
	resource := input.Resource // ElasticPool
	parent := input.Parent     // Server

	fullParent := input.Parent
	if input.FullParent != "" {
		fullParent = input.FullParent
	}

	idPackage, idType, idName := parseIDStructure(input.IDStructure)

	parentLower := strings.ToLower(parent)
	serviceLower := strings.ToLower(service)
	resourceLower := strings.ToLower(resource)
	useResourceGroup := false
	if parentLower == "resourcegroup" || parentLower == "" {
		useResourceGroup = true
	}
	terraformName := fmt.Sprintf(
		"azurerm_%s_%s",
		serviceLower,
		resourceLower,
	)
	resourceName := fmt.Sprintf(
		"%s_%s",
		serviceLower,
		resourceLower,
	)
	if input.TerraformName != "" {
		resourceName = input.TerraformName
		terraformName = fmt.Sprintf(
			"azurerm_%s",
			input.TerraformName,
		)
	}

	parentTerraformName := parent

	if input.ParentTerraformName != "" {
		parentTerraformName = input.ParentTerraformName
	}

	listResourceStruct := fmt.Sprintf("%sListResource", resource)
	resourceStruct := fmt.Sprintf("%sResource", resource)
	modelStruct := fmt.Sprintf("%sListModel", resource)

	childDisplayName := resource

	if input.IncludeServiceInName {
		resourceStruct = fmt.Sprintf("%s%sResource", service, resource)
		listResourceStruct = fmt.Sprintf("%s%sListResource", service, resource)
		modelStruct = fmt.Sprintf("%s%sListModel", service, resource)

		childDisplayName = fmt.Sprintf("%s %s", service, resource)
	}

	fmt.Printf("return resource%s%sFlatten(d, id, resp.Model) \n}\n\nfunc resource%s%sFlatten(d *pluginsdk.ResourceData, id *%s.%sId, model *%s.%s) error {\n\n", service, resource, service, resource, idPackage, idType, idPackage, idType)

	return &ListResourceData{
		// Core
		ServiceName:  service,
		ResourceName: resource,
		PackageName:  serviceLower,

		ResourceLower:    resourceLower,
		FullParent:       fullParent,
		Parent:           parent,
		ParentLower:      parentLower,
		UseResourceGroup: useResourceGroup,

		IdPackage: idPackage,
		IdType:    idType,
		IdName:    idName,
		// Terraform
		TerraformResourceName:       terraformName,
		ParentTerraformResourceName: parentTerraformName,

		// Go identifiers
		ResourceStruct:     resourceStruct,
		ListResourceStruct: listResourceStruct,
		ModelStruct:        modelStruct,

		// Display
		ChildDisplayName: childDisplayName,

		// Files
		OutputFile:     fmt.Sprintf("/%s_resource_list.go", resourceName),
		OutputTestFile: fmt.Sprintf("/%s_resource_list_test.go", resourceName),
	}
}

func parseIDStructure(idStructure string) (string, string, string) {
	parts := strings.Split(idStructure, ".")
	name := strings.TrimSuffix(parts[1], "Id") // Remove "Id" suffix for type name
	return parts[0], parts[1], name
}
