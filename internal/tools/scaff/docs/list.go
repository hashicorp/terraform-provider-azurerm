// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package docs generates list-resource documentation (Markdown) from the
// generated *_resource_list.go files produced by the scaff tool. It is invoked
// automatically whenever a list resource is generated or updated so the
// documentation under website/docs/list-resources stays in sync with the code.
package docs

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/list_document.tmpl
var templatesFS embed.FS

type TemplateData struct {
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

// GenerateListDocs walks root looking for generated *_resource_list.go files and
// generates the corresponding list-resource documentation for each, returning
// the paths of the files written. An existing document is left untouched unless
// overwrite is true.
func GenerateListDocs(root string, overwrite bool) ([]string, error) {
	var written []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_resource_list.go") {
			return nil
		}
		out, ok, err := GenerateListDoc(path, overwrite)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if ok {
			written = append(written, out)
		}
		return nil
	})
	return written, err
}

// GenerateListDoc parses a single generated *_resource_list.go file and writes
// the corresponding list-resource documentation Markdown to
// website/docs/list-resources, returning the output path and whether the file
// was written. An existing document is preserved (written == false) unless
// overwrite is true, so curated documentation is not clobbered.
func GenerateListDoc(filename string, overwrite bool) (string, bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return "", false, err
	}
	var hasListSchema bool
	var (
		resourceName string
		attributes   []Attribute
	)

	ast.Inspect(node, func(n ast.Node) bool {
		// Detect ListResourceConfigSchema method
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "ListResourceConfigSchema" {
			hasListSchema = true
		}

		// Metadata → TypeName (unchanged)
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "Metadata" {
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					val := strings.Trim(lit.Value, "`\"")
					if strings.HasPrefix(val, "azurerm_") {
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
		return "", false, fmt.Errorf("resource type name not found in %s", filename)
	}

	if !hasListSchema || len(attributes) == 0 {
		attributes = defaultListAttributes
	}

	outDir, err := markdownPathFromGo(filename, "", true)
	if err != nil {
		return "", false, err
	}
	base := strings.TrimSuffix(filepath.Base(filename), "_resource_list.go")
	outPath := filepath.Join(outDir, base+".html.markdown")

	// Preserve curated documentation: only replace an existing file when the
	// caller explicitly opts in to overwriting generated artifacts.
	if !overwrite {
		if _, err := os.Stat(outPath); err == nil {
			return outPath, false, nil
		}
	}

	md, err := renderMarkdown(resourceName, attributes, filename)
	if err != nil {
		return "", false, err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return "", false, err
	}
	if err := os.WriteFile(outPath, []byte(md), 0o644); err != nil {
		return "", false, err
	}

	return outPath, true, nil
}

func renderMarkdown(resource string, attrs []Attribute, filename string) (string, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/list_document.tmpl")
	if err != nil {
		return "", err
	}

	data := buildTemplateData(resource, attrs, filename)

	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}

func buildTemplateData(resource string, attrs []Attribute, filename string) TemplateData {

	data := TemplateData{
		Resource:      resource,
		Section:       mapSectionToActualSection(getSectionFromResourceTitle(friendlyTitle(resource))),
		FriendlyTitle: friendlyTitle(resource),
	}

	for _, a := range attrs {
		// Build examples

		exampleNameBit := ""

		if a.Name == "resource_group_name" {
			exampleNameBit = "rg"
		}

		switch {
		case a.Name == "subscription_id":
			data.Examples = append(data.Examples, Example{
				Heading: fmt.Sprintf("List all %ss in the subscription", data.FriendlyTitle),
			})

		case isNameAttribute(a.Name):
			data.Examples = append(data.Examples, Example{
				Heading:               fmt.Sprintf("List all %ss in a %s", data.FriendlyTitle, nameFriendlyName(a.Name)),
				AttributeName:         a.Name,
				AttributeExampleValue: fmt.Sprintf(`"example-%s"`, exampleNameBit),
			})

		default:
			idExamplePath, err := markdownPathFromGo(filename, idRemoveID(checkAddSectionToName(getSectionFromResourceName(resource), a.Name)), false)
			if err != nil {
				continue
			}
			idExample := getExampleValueForIDAttribute(idExamplePath)
			if idExample == "" {
				// The resource's own documentation may not exist yet (e.g. when a
				// list resource is generated ahead of its docs); fall back to a
				// placeholder so the example stays valid HCL.
				idExample = `"..."`
			}

			data.Examples = append(data.Examples, Example{
				Heading:               fmt.Sprintf("List %ss in a %s", data.FriendlyTitle, idFriendlyName(checkAddSectionToName(getSectionFromResourceName(resource), a.Name))),
				AttributeName:         a.Name,
				AttributeExampleValue: idExample,
			})
		}

		// Build arguments
		req := "Required"
		if a.Optional {
			req = "Optional"
		}

		extraBit := ""

		if a.Name == "subscription_id" {
			extraBit = " Defaults to the value specified in the Provider Configuration."
		}

		if isIDAttribute(a.Name) {
			data.Arguments = append(data.Arguments, Argument{
				Name:        a.Name,
				Requirement: req,
				Description: fmt.Sprintf("The ID of the %s to query.%s", idFriendlyName(checkAddSectionToName(getSectionFromResourceName(resource), a.Name)), extraBit),
			})
		}

		if isNameAttribute(a.Name) {
			data.Arguments = append(data.Arguments, Argument{
				Name:        a.Name,
				Requirement: req,
				Description: fmt.Sprintf("The name of the %s to query.", nameFriendlyName(a.Name)),
			})
		}
	}

	return data
}

func isNameAttribute(name string) bool {
	return strings.HasSuffix(name, "_name")
}

func nameFriendlyName(name string) string {
	base := strings.TrimSuffix(name, "_name")
	parts := strings.Split(base, "_")

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}

	return strings.Join(parts, " ")
}

func isIDAttribute(name string) bool {
	return strings.HasSuffix(name, "_id")
}

func idRemoveID(name string) string {
	return strings.TrimSuffix(name, "_id")
}

func idFriendlyName(name string) string {
	base := strings.TrimSuffix(name, "_id")
	parts := strings.Split(base, "_")

	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}

	return strings.Join(parts, " ")
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

	sectionsCheck := []string{"mssql"}
	for _, s := range sectionsCheck {
		if s == section {
			return s + "_" + name
		}
	}

	return name
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
	isAbs := filepath.IsAbs(goPath)

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
		return "", fmt.Errorf("could not determine service from path %q", goPath)
	}

	// filepath.Join drops a leading separator, so restore it when the input was
	// an absolute path to keep the derived documentation path absolute too. This
	// lets the generator accept both the repo-relative paths used by the scaff
	// commands and absolute paths.
	if isAbs && moduleRoot != "" && !filepath.IsAbs(moduleRoot) {
		moduleRoot = string(filepath.Separator) + moduleRoot
	}

	if isOutput {
		return filepath.Join(moduleRoot, "website", "docs", "list-resources"), nil
	}
	return filepath.Join(moduleRoot, "website", "docs", "r", attributeName+".html.markdown"), nil
}

func getExampleValueForIDAttribute(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
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
