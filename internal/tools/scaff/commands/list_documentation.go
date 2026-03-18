package commands

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/templatehelpers"
	"github.com/mitchellh/cli"
)

type ListDocumentationCommand struct {
	Ui cli.Ui
}

type ListDocumentationConfig struct {
	Path             string
	SubCategory      string
	AddSectionToName bool
}

var _ cli.Command = &ListDocumentationCommand{}

type ListDocumentationData struct {
	Resource      string
	SubCategory   string
	FriendlyTitle string
	Examples      []Example
	Arguments     []Argument
}

type Example struct {
	Heading               string
	AttributeName         string
	AttributeExampleValue string
}

type Argument struct {
	Name        string
	Requirement string
	Description string
}

type Attribute struct {
	Name     string
	Optional bool
}

var defaultListAttributes = []Attribute{
	{
		Name:     "subscription_id",
		Optional: true,
	},
	{
		Name:     "resource_group_name",
		Optional: true,
	},
}

func (c ListDocumentationCommand) Run(args []string) int {
	config := &ListDocumentationConfig{}

	if err := config.parseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	err := filepath.WalkDir(config.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), "_resource_list.go") {
			if err := c.processFile(path, config); err != nil {
				c.Ui.Error(fmt.Sprintf("❌ %s: %v", path, err))
			} else {
				c.Ui.Info(fmt.Sprintf("✅ processed %s", path))
			}
		}

		return nil
	})
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (config *ListDocumentationConfig) parseArgs(args []string) error {
	argSet := flag.NewFlagSet("list-documentation", flag.ContinueOnError)

	argSet.StringVar(&config.Path, "path", "", "(Required) Path to file or directory to scan")
	argSet.StringVar(&config.SubCategory, "subcategory", "", "(Optional) Subcategory to override section mapping")
	argSet.BoolVar(&config.AddSectionToName, "addsectiontoname", false, "(Optional) Prepend section name to attribute names")

	if err := argSet.Parse(args); err != nil {
		return err
	}

	if config.Path == "" {
		return errors.New("path is required")
	}

	return nil
}

func (c ListDocumentationCommand) Synopsis() string {
	return "Generate list resource documentation from Go source files"
}

func (c ListDocumentationCommand) Help() string {
	return `
Usage: scaff list-documentation [options]

  Generates documentation for list resources by scanning Go source files
  for files ending with '_resource_list.go'.

Options:
  -path=<path>              (Required) Path to file or directory to scan (required)
  -subcategory=<name>       (Optional) Override the subcategory/section (e.g., "Network", "Database")
  -addsectiontoname         (Optional) Prepend section name to attribute names (boolean flag)

Examples:
  # Basic usage (backward compatible)
  scaff list-documentation internal/services/network

  # Using flags
  scaff list-documentation -path=internal/services/network

  # With subcategory override
  scaff list-documentation -path=internal/services/mssql -subcategory=Database

  # With section name prefixing enabled
  scaff list-documentation -path=internal/services/mssql -addsectiontoname
`
}

func (c ListDocumentationCommand) processFile(filename string, config *ListDocumentationConfig) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return err
	}
	var hasListSchema bool
	var (
		resourceName    string
		attributes      []Attribute
		varDeclarations = make(map[string]string) // Map variable names to their string values
	)

	// First pass: collect variable declarations from current file
	ast.Inspect(node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range vspec.Names {
						if i < len(vspec.Values) {
							if lit, ok := vspec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								val := strings.Trim(lit.Value, "`\"")
								varDeclarations[name.Name] = val
							}
						}
					}
				}
			}
		}
		return true
	})

	// Also check the corresponding resource file (without _list suffix) for variable declarations
	if strings.HasSuffix(filename, "_resource_list.go") {
		resourceFilename := strings.Replace(filename, "_resource_list.go", "_resource.go", 1)
		if resourceNode, err := parser.ParseFile(fset, resourceFilename, nil, 0); err == nil {
			ast.Inspect(resourceNode, func(n ast.Node) bool {
				if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
					for _, spec := range decl.Specs {
						if vspec, ok := spec.(*ast.ValueSpec); ok {
							for i, name := range vspec.Names {
								if i < len(vspec.Values) {
									if lit, ok := vspec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
										val := strings.Trim(lit.Value, "`\"")
										varDeclarations[name.Name] = val
									}
								}
							}
						}
					}
				}
				return true
			})
		}
	}

	// Second pass: find resource name and attributes
	ast.Inspect(node, func(n ast.Node) bool {
		// Detect ListResourceConfigSchema method
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "ListResourceConfigSchema" {
			hasListSchema = true
		}

		// Metadata → TypeName - handle both string literals and variable references
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "Metadata" {
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				// Handle direct string literal: response.TypeName = "azurerm_..."
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					val := strings.Trim(lit.Value, "`\"")
					if strings.HasPrefix(val, "azurerm_") {
						resourceName = val
					}
				}
				// Handle variable reference: response.TypeName = variableName
				if ident, ok := n.(*ast.Ident); ok {
					if val, exists := varDeclarations[ident.Name]; exists && strings.HasPrefix(val, "azurerm_") {
						resourceName = val
					}
				}
				return true
			})
		}

		// Schema attributes
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		for _, elt := range cl.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			key, ok := kv.Key.(*ast.BasicLit)
			if !ok || key.Kind != token.STRING {
				continue
			}

			attr := Attribute{
				Name: strings.Trim(key.Value, "`\""),
			}

			if lit, ok := kv.Value.(*ast.CompositeLit); ok {
				for _, e := range lit.Elts {
					if kv, ok := e.(*ast.KeyValueExpr); ok {
						if id, ok := kv.Key.(*ast.Ident); ok && id.Name == "Optional" {
							if v, ok := kv.Value.(*ast.Ident); ok && v.Name == "true" {
								attr.Optional = true
							}
						}
					}
				}
			}

			attributes = append(attributes, attr)
		}

		return true
	})

	if resourceName == "" {
		return fmt.Errorf("resource type name not found")
	}

	if !hasListSchema || len(attributes) == 0 {
		attributes = defaultListAttributes
	}

	md, err := config.renderMarkdown(resourceName, attributes, filename)
	if err != nil {
		return err
	}

	docsDirectory, err := markdownPathFromGo(filename, "", true)
	if err != nil {
		return err
	}

	baseFileName := strings.TrimSuffix(filepath.Base(filename), "_resource_list.go")
	outputPath := filepath.Join(docsDirectory, baseFileName+".html.markdown")

	if err := os.WriteFile("/"+outputPath, []byte(md), 0o644); err != nil {
		return err
	}

	c.Ui.Info(fmt.Sprintf("✅ generated %s", outputPath))
	return nil
}

func (config *ListDocumentationConfig) renderMarkdown(resource string, attrs []Attribute, filename string) (string, error) {
	tmpl := template.Must(template.New("list_document.html.markdown.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/list_document.html.markdown.gotpl"))

	data, err := config.buildTemplateData(resource, attrs, filename)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (config *ListDocumentationConfig) buildExampleForAttribute(a Attribute, resourceTitle string, section string, filename string) (Example, error) {
	exampleNameSuffix := ""
	if a.Name == "resource_group_name" {
		exampleNameSuffix = "-rg"
	}
	attributeName := config.checkAddSectionToName(section, a.Name)

	switch {
	case a.Name == "subscription_id":
		return Example{
			Heading: fmt.Sprintf("List all %ss in the subscription", resourceTitle),
		}, nil

	case isNameAttribute(a.Name):
		return Example{
			Heading:               fmt.Sprintf("List all %ss in a %s", resourceTitle, nameTitle(a.Name)),
			AttributeName:         a.Name,
			AttributeExampleValue: fmt.Sprintf(`"example%s"`, exampleNameSuffix),
		}, nil

	default:
		idExamplePath, err := markdownPathFromGo(filename, idRemoveID(attributeName), false)
		if err != nil {
			return Example{}, fmt.Errorf("failed to derive markdown path for example value for attribute `%s`: %w", a.Name, err)
		}
		idExample := getExampleValueForIDAttribute(idExamplePath)

		return Example{
			Heading:               fmt.Sprintf("List %ss in a %s", resourceTitle, idTitle(attributeName)),
			AttributeName:         a.Name,
			AttributeExampleValue: idExample,
		}, nil
	}
}

func (config *ListDocumentationConfig) buildArgumentForAttribute(a Attribute, sectionFromName string) (Argument, error) {
	req := "Required"
	if a.Optional {
		req = "Optional"
	}

	extraDescription := ""
	if a.Name == "subscription_id" {
		extraDescription = " Defaults to the value specified in the Provider Configuration."
	}

	friendlyName := config.checkAddSectionToName(sectionFromName, a.Name)

	if isIDAttribute(a.Name) {
		return Argument{
			Name:        a.Name,
			Requirement: req,
			Description: fmt.Sprintf("The ID of the %s to query.%s", idTitle(friendlyName), extraDescription),
		}, nil
	}

	if isNameAttribute(a.Name) {
		return Argument{
			Name:        a.Name,
			Requirement: req,
			Description: fmt.Sprintf("The name of the %s to query.", nameTitle(friendlyName)),
		}, nil
	}

	return Argument{}, fmt.Errorf("%s is not an ID or name attribute", a.Name)
}

func (config *ListDocumentationConfig) buildTemplateData(resource string, attrs []Attribute, filename string) (ListDocumentationData, error) {
	resourceTitle := resourceTitle(resource)
	section := getSection(resource)

	subCategory := section
	if config.SubCategory != "" {
		subCategory = config.SubCategory
	}

	data := ListDocumentationData{
		Resource:      resource,
		SubCategory:   subCategory,
		FriendlyTitle: resourceTitle,
	}

	for _, a := range attrs {
		example, err := config.buildExampleForAttribute(a, resourceTitle, section, filename)
		if err != nil {
			return ListDocumentationData{}, err
		}
		data.Examples = append(data.Examples, example)

		arg, err := config.buildArgumentForAttribute(a, section)
		if err != nil {
			return ListDocumentationData{}, err
		}
		data.Arguments = append(data.Arguments, arg)
	}

	return data, nil
}

func isNameAttribute(name string) bool {
	return strings.HasSuffix(name, "_name")
}

func isIDAttribute(name string) bool {
	return strings.HasSuffix(name, "_id")
}

func idRemoveID(name string) string {
	return strings.TrimSuffix(name, "_id")
}

func toTitle(name, suffix string) string {
	base := strings.TrimSuffix(strings.TrimPrefix(name, "azurerm_"), suffix)
	parts := strings.Split(base, "_")

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}

	return strings.Join(parts, " ")
}

func nameTitle(name string) string {
	return toTitle(name, "_name")
}

func idTitle(name string) string {
	return toTitle(name, "_id")
}

func resourceTitle(name string) string {
	return toTitle(name, "")
}

func (config *ListDocumentationConfig) checkAddSectionToName(section string, name string) string {
	if name == "subscription_id" || name == "resource_group_name" {
		return name
	}

	if config.AddSectionToName {
		return section + "_" + name
	}

	return name
}

func getSection(name string) string {
	parts := strings.Split(name, "_")

	return parts[1]
}

func markdownPathFromGo(goPath string, attributeName string, isOutput bool) (string, error) {
	// Normalize path
	goPath = filepath.Clean(goPath)

	parts := strings.Split(goPath, string(filepath.Separator))

	// Find "internal/services"
	var moduleRoot string
	var service string

	for i := 0; i < len(parts)-2; i++ {
		if parts[i] == "internal" && parts[i+1] == "services" {
			service = parts[i+2]
			moduleRoot = filepath.Join(parts[:i]...)
			break
		}
	}

	if service == "" {
		return "", fmt.Errorf("could not determine service from path")
	}

	if isOutput {
		outputPath := filepath.Join(moduleRoot, "website", "docs", "list-resources")
		return outputPath, nil
	} else {
		mdPath := filepath.Join(moduleRoot, "website", "docs", "r", attributeName+".html.markdown")
		return mdPath, nil
	}
}

func getExampleValueForIDAttribute(filepath string) string {
	file, err := os.Open("/" + filepath)
	if err != nil {
		return fmt.Sprintf("failed to open %s", filepath)
	}
	defer file.Close()

	re := regexp.MustCompile(`^\s*terraform import\s+[^\s]+\s+(.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := re.FindStringSubmatch(line); matches != nil {
			return fmt.Sprintf("%q", matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}
