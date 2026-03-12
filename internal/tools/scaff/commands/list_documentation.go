package commands

import (
	"bufio"
	"bytes"
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

var _ cli.Command = &ListDocumentationCommand{}

type ListDocumentationData struct {
	Resource      string
	Section       string
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
	if len(args) < 1 {
		c.Ui.Error("usage: scaff list-documentation <directory>")
		return 1
	}

	root := args[0]

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), "_resource_list.go") {
			if err := c.processFile(path); err != nil {
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

func (c ListDocumentationCommand) Synopsis() string {
	return "Generate list resource documentation from Go source files"
}

func (c ListDocumentationCommand) Help() string {
	return `
Usage: scaff list-documentation <file>

  Generates documentation for list resources by scanning Go source files
  in the specified directory for files ending with '_resource_list.go'.

Parameters:
  <directory> (Required) The file or directory to scan for list resource files.

Example:
  scaff list-documentation internal/services/network
`
}

func (c ListDocumentationCommand) processFile(filename string) error {
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

	md, err := renderMarkdown(resourceName, attributes, filename)
	if err != nil {
		return err
	}

	base := strings.TrimSuffix(filepath.Base(filename), "_resource_list.go")

	out, err := markdownPathFromGo(filename, "", true)
	if err != nil {
		return err
	}
	out = filepath.Join(out, base+".html.markdown")

	if err := os.WriteFile("/"+out, []byte(md), 0644); err != nil {
		return err
	}

	c.Ui.Info(fmt.Sprintf("✅ generated %s", out))
	return nil
}

func renderMarkdown(resource string, attrs []Attribute, filename string) (string, error) {
	tmpl := template.Must(template.New("list_document.html.markdown.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/list_document.html.markdown.gotpl"))

	data, err := buildTemplateData(resource, attrs, filename)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func buildExampleForAttribute(a Attribute, friendlyTitleValue, sectionFromName, filename string) (Example, error) {
	exampleNameBit := ""
	if a.Name == "resource_group_name" {
		exampleNameBit = "-rg"
	}

	switch {
	case a.Name == "subscription_id":
		return Example{
			Heading: fmt.Sprintf("List all %ss in the subscription", friendlyTitleValue),
		}, nil

	case isNameAttribute(a.Name):
		return Example{
			Heading:               fmt.Sprintf("List all %ss in a %s", friendlyTitleValue, nameFriendlyName(a.Name)),
			AttributeName:         a.Name,
			AttributeExampleValue: fmt.Sprintf(`"example%s"`, exampleNameBit),
		}, nil

	default:
		idExamplePath, err := markdownPathFromGo(filename, idRemoveID(checkAddSectionToName(sectionFromName, a.Name)), false)
		if err != nil {
			return Example{}, fmt.Errorf("failed to derive markdown path for example value for attribute `%s`: %w", a.Name, err)
		}
		idExample := getExampleValueForIDAttribute(idExamplePath)

		return Example{
			Heading:               fmt.Sprintf("List %ss in a %s", friendlyTitleValue, idFriendlyName(checkAddSectionToName(sectionFromName, a.Name))),
			AttributeName:         a.Name,
			AttributeExampleValue: idExample,
		}, nil
	}
}

func buildArgumentForAttribute(a Attribute, sectionFromName string) (Argument, bool) {
	req := "Required"
	if a.Optional {
		req = "Optional"
	}

	extraBit := ""
	if a.Name == "subscription_id" {
		extraBit = " Defaults to the value specified in the Provider Configuration."
	}

	if isIDAttribute(a.Name) {
		return Argument{
			Name:        a.Name,
			Requirement: req,
			Description: fmt.Sprintf("The ID of the %s to query.%s", idFriendlyName(checkAddSectionToName(sectionFromName, a.Name)), extraBit),
		}, true
	}

	if isNameAttribute(a.Name) {
		return Argument{
			Name:        a.Name,
			Requirement: req,
			Description: fmt.Sprintf("The name of the %s to query.", nameFriendlyName(a.Name)),
		}, true
	}

	return Argument{}, false
}

func buildTemplateData(resource string, attrs []Attribute, filename string) (ListDocumentationData, error) {
	friendlyTitleValue := friendlyTitle(resource)
	sectionFromTitle := getSectionFromResourceTitle(friendlyTitleValue)
	mappedSection := mapSectionToActualSection(sectionFromTitle)
	sectionFromName := getSectionFromResourceName(resource)

	data := ListDocumentationData{
		Resource:      resource,
		Section:       mappedSection,
		FriendlyTitle: friendlyTitleValue,
	}

	for _, a := range attrs {
		example, err := buildExampleForAttribute(a, friendlyTitleValue, sectionFromName, filename)
		if err != nil {
			return ListDocumentationData{}, fmt.Errorf("failed to build example for attribute %s: %w", a.Name, err)
		}
		data.Examples = append(data.Examples, example)

		if arg, shouldAdd := buildArgumentForAttribute(a, sectionFromName); shouldAdd {
			data.Arguments = append(data.Arguments, arg)
		}
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

// toFriendlyName converts snake_case to Title Case, removing the specified suffix
func toFriendlyName(name, suffix string) string {
	base := strings.TrimSuffix(name, suffix)
	parts := strings.Split(base, "_")

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}

	return strings.Join(parts, " ")
}

func nameFriendlyName(name string) string {
	return toFriendlyName(name, "_name")
}

func idFriendlyName(name string) string {
	return toFriendlyName(name, "_id")
}

func friendlyTitle(name string) string {

	parts := strings.Split(strings.TrimPrefix(name, "azurerm_"), "_")

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, " ")
}

func checkAddSectionToName(section string, name string) string {
	if name == "subscription_id" || name == "resource_group_name" {
		return name
	}

	sectionsCheck := []string{"mysql"}
	for _, s := range sectionsCheck {
		if s == section {
			return s + "_" + name
		}
	}

	return name
}

func removeSectionFromResourceName(section string, resource string) string {
	if strings.HasPrefix(resource, "azurerm_"+section+"_") {
		resource = strings.TrimPrefix(resource, "azurerm_"+section+"_")
	}
	resource = strings.TrimPrefix(resource, "azurerm_")
	return resource
}

func getSectionFromResourceTitle(title string) string {
	parts := strings.Split(title, " ")

	return parts[0]
}

func getSectionFromResourceName(name string) string {
	parts := strings.Split(name, "_")

	return parts[1]
}

func mapSectionToActualSection(section string) string {

	var sectionCategoryMap = map[string]string{
		"Firewall":    "Network",
		"Mssql":       "Database",
		"Mysql":       "Database",
		"Application": "Network",
		"Ip":          "Network",
		"Web":         "Network",
		"Route":       "Network",
		"Public":      "Network",
		"Private":     "Network",
		"Nat":         "Network",
		"Virtual":     "Network",
	}

	if category, ok := sectionCategoryMap[section]; ok {
		return category
	}
	return section
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
		return "failed to open"
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
