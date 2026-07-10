package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/cli"
)

type DocumentCommand struct {
	Ui cli.Ui
}

var _ cli.Command = DocumentCommand{}

type DocumentData struct {
	Name                  string
	SnakeName             string
	ProviderName          string
	ProviderCanonicalName string
	ServicePackage        string
	DocType               string
	Resource              helpers.Resource
	IDExample             string
	Examples              []string
	ResourceData          string
	SchemaAPIURL          string
}

func (d *DocumentData) ParseArgs(args []string) (errors []error) {
	docSet := flag.NewFlagSet("document", flag.ExitOnError)
	docSet.StringVar(&d.Name, "name", "", "The name of the resource")
	docSet.StringVar(&d.ServicePackage, "servicepackage", "", "The name of the Service Package the resource or data source belongs to")
	docSet.StringVar(&d.DocType, "type", "", "The type of item to document, one of `resource` or `data-source`")
	docSet.StringVar(&d.IDExample, "id", "", "An example of the ID this resource has when created (only valid for `-type=resource`)")
	docSet.StringVar(&d.SchemaAPIURL, "schemaapiurl", config.SchemaAPIURL, "The URL of the Provider's JSON API (defaults to http://localhost:8080)")
	err := docSet.Parse(args)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.Name == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
		return errors
	}
	if d.DocType == "" {
		errors = append(errors, fmt.Errorf("required option `-type` missing\n"))
		return errors
	}

	if strings.EqualFold(d.DocType, "resource") && d.IDExample == "" {
		errors = append(errors, fmt.Errorf("`-id` required when `-type=resource\n`"))
	}

	if strings.EqualFold(d.DocType, "datasource") {
		d.DocType = "data-source"
	}

	if config.ProviderCanonicalName != "" {
		d.ProviderCanonicalName = config.ProviderCanonicalName
	} else {
		d.ProviderCanonicalName = d.Name
	}

	return errors
}

func (c DocumentCommand) Run(args []string) int {
	data := &DocumentData{}

	if len(args) == 0 {
		fmt.Print(c.Help())
		return 1
	}

	if err := data.ParseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generate(); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to generate: %+v", err))
		return 1
	}

	return 0
}

func (c DocumentCommand) Help() string {
	return "TODO"
}

func (c DocumentCommand) Synopsis() string {
	return "generates documentation from a resource"
}

func (d *DocumentData) generate() error {
	if d.ProviderName == "" {
		providerName, err := helpers.ProviderName()
		if err != nil {
			return err
		}
		d.ProviderName = *providerName
	}
	if d.SnakeName == "" {
		d.SnakeName = strcase.ToSnake(fmt.Sprintf("%s_%s", d.ProviderName, d.Name))
	}

	resource, err := helpers.GetSchema(d.SchemaAPIURL, d.DocType, d.SnakeName)
	if err != nil {
		return err
	}
	d.Resource = *resource

	tpl := template.Must(template.New("").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/document-*.gotpl"))

	outputPath := config.DocsPath
	if d.DocType == "data-source" {
		outputPath = fmt.Sprintf("%s/%s", outputPath, config.DataSourceDocsDirname)
	} else {
		outputPath = fmt.Sprintf("%s/%s", outputPath, config.ResourceDocsDirname)
	}

	providerPrefix := fmt.Sprintf("%s_", d.ProviderName)
	outputPath = fmt.Sprintf("%s/%s.md", outputPath, strings.TrimPrefix(d.SnakeName, providerPrefix))

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("[DEBUG] failed opening output documentation file for writing: %+v", err.Error())
		return err
	}

	defer f.Close()

	return d.renderDoc(f, tpl)
}

func (d *DocumentData) renderDoc(f *os.File, tpl *template.Template) error {
	if d.DocType == string(helpers.DocTypeDataSource) {
		switch config.DocsVersion {
		case helpers.DocsVersionRegistry:
			return d.renderDataSourceDocRegistry(f, tpl)
		default:
			return d.renderDataSourceDocLegacy(f, tpl)
		}
	}

	switch config.DocsVersion {
	case helpers.DocsVersionRegistry:
		return d.renderResourceDocRegistry(f, tpl)
	default:
		return d.renderResourceDocLegacy(f, tpl)
	}
}

func (d *DocumentData) renderResourceDocLegacy(f *os.File, tpl *template.Template) error {
	if err := tpl.ExecuteTemplate(f, "document-frontmatter-legacy-resource.gotpl", d); err != nil {
		return err
	}
	if err := tpl.ExecuteTemplate(f, "document-header-resource-legacy.gotpl", d); err != nil {
		return err
	}
	if err := tpl.ExecuteTemplate(f, "document-usage-example.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-required.gotpl", d); err != nil {
		return err
	}

	// only process the Optional params if at least one top level item is Optional
	for _, v := range d.Resource.Schema {
		if v.Optional {
			if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-optional.gotpl", d); err != nil {
				return err
			}
			break
		}
	}

	for k, v := range d.Resource.Schema {
		if data := v.Elem; data != nil {
			if s, ok := data.(map[string]interface{}); ok {
				if _, ok := s["schema"].(map[string]interface{}); ok {
					if err := renderNestedSchemaBlock(k, s, f, tpl); err != nil {
						return err
					}
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-attributes-legacy.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-footer-resource.gotpl", d); err != nil {
		return err
	}

	return nil
}

func (d *DocumentData) renderResourceDocRegistry(f *os.File, tpl *template.Template) error {
	if err := tpl.ExecuteTemplate(f, "document-frontmatter-registry.gotpl", d); err != nil {
		return err
	}
	if err := tpl.ExecuteTemplate(f, "document-header-resource-registry.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-usage-example.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-required-registry.gotpl", d); err != nil {
		return err
	}

	// only process the Optional params if at least one top level item is Optional
	for _, v := range d.Resource.Schema {
		if v.Optional {
			if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-optional-registry.gotpl", d); err != nil {
				return err
			}
			break
		}
	}

	for k, v := range d.Resource.Schema {
		if data := v.Elem; data != nil {
			if s, ok := data.(map[string]interface{}); ok {
				if _, ok := s["schema"].(map[string]interface{}); ok {
					if err := renderNestedSchemaBlock(k, s, f, tpl); err != nil {
						return err
					}
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-attributes-registry.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-footer-resource.gotpl", d); err != nil {
		return err
	}

	return nil
}

func (d *DocumentData) renderDataSourceDocLegacy(f *os.File, tpl *template.Template) error {
	if err := tpl.ExecuteTemplate(f, "document-frontmatter-legacy-datasource.gotpl", d); err != nil {
		return err
	}
	if err := tpl.ExecuteTemplate(f, "document-header-datasource-legacy.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-usage-example.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-required.gotpl", d); err != nil {
		return err
	}

	// only process the Optional params if at least one top level item is Optional
	for _, v := range d.Resource.Schema {
		if v.Optional {
			if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-optional.gotpl", d); err != nil {
				return err
			}
			break
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-attributes-legacy.gotpl", d); err != nil {
		return err
	}

	for k, v := range d.Resource.Schema {
		if data := v.Elem; data != nil {
			if s, ok := data.(map[string]interface{}); ok {
				if _, ok := s["schema"].(map[string]interface{}); ok {
					if err := renderNestedSchemaBlock(k, s, f, tpl); err != nil {
						return err
					}
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-footer-datasource.gotpl", d); err != nil {
		return err
	}

	return nil
}

func (d *DocumentData) renderDataSourceDocRegistry(f *os.File, tpl *template.Template) error {
	if err := tpl.ExecuteTemplate(f, "document-frontmatter-registry.gotpl", d); err != nil {
		return err
	}
	if err := tpl.ExecuteTemplate(f, "document-header-datasource-registry.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-usage-example.gotpl", d); err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-required-registry.gotpl", d); err != nil {
		return err
	}

	// only process the Optional params if at least one top level item is Optional
	for _, v := range d.Resource.Schema {
		if v.Optional {
			if err := tpl.ExecuteTemplate(f, "document-schema-toplevel-optional-registry.gotpl", d); err != nil {
				return err
			}
			break
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-schema-attributes-registry.gotpl", d); err != nil {
		return err
	}

	for k, v := range d.Resource.Schema {
		if data := v.Elem; data != nil {
			if s, ok := data.(map[string]interface{}); ok {
				if _, ok := s["schema"].(map[string]interface{}); ok {
					if err := renderNestedSchemaBlock(k, s, f, tpl); err != nil {
						return err
					}
				}
			}
		}
	}

	if err := tpl.ExecuteTemplate(f, "document-footer-datasource.gotpl", d); err != nil {
		return err
	}

	return nil
}

func renderNestedSchemaBlock(name string, props map[string]interface{}, f *os.File, tpl *template.Template) error {
	if err := tpl.ExecuteTemplate(f, "document-schema-block-header.gotpl", name); err != nil {
		return err
	}
	if config.DocsVersion == helpers.DocsVersionRegistry {
		if err := tpl.ExecuteTemplate(f, "document-schema-block-registry.gotpl", props); err != nil {
			return err
		}
	} else {
		if err := tpl.ExecuteTemplate(f, "document-schema-block-legacy.gotpl", props); err != nil {
			return err
		}
	}

	if schema, ok := props["schema"].(map[string]interface{}); ok {
		for k, v := range schema {
			if d, ok := v.(map[string]interface{}); ok && d["elem"] != nil {
				if data := helpers.FlattenMapToSchema(d).Elem; data != nil {
					if s, ok := data.(map[string]interface{}); ok {
						if _, ok := s["schema"].(map[string]interface{}); ok {
							if err := renderNestedSchemaBlock(k, s, f, tpl); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}
