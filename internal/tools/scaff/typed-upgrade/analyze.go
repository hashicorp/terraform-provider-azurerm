// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
	list_upgrade "github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/list-upgrade"
)

// Info holds the full typed-upgrade analysis of an untyped resource.
type Info struct {
	// Base metadata derived from list_upgrade.Analyze.
	Path          string
	Package       string
	BaseName      string // "MapsCreator" from "resourceMapsCreator"
	GoName        string // same as BaseName for the struct root
	StructName    string // GoName + "Resource"
	ModelName     string // GoName + "ResourceModel"
	TerraformType string // derived, e.g. "azurerm_maps_creator"

	// SDK + ID metadata (from list-upgrade analysis).
	SDKPackage     string
	SDKImportPath  string
	IDPackage      string
	IDImportPath   string
	IDTypeName     string
	IDBase         string
	IDParseFunc    string
	IDValidateFunc string
	ServiceName    string
	ClientField    string

	// CRUD function names.
	CreateFunc string
	ReadFunc   string
	UpdateFunc string
	DeleteFunc string

	// Flags.
	HasUpdate        bool
	HasCustomizeDiff bool
	HasSchemaVersion bool
	SchemaVersion    int
	DeprecationMsg   string

	// Schema fields (all top-level entries in the Schema map).
	Fields []*SchemaField

	// Timeouts per operation (raw source expressions, e.g. "30 * time.Minute").
	CreateTimeout string
	ReadTimeout   string
	UpdateTimeout string
	DeleteTimeout string

	// CRUD function bodies ready for wrapping (raw source of the function body,
	// with outer braces stripped, and safe text-level migrations applied).
	CreateBody string
	ReadBody   string
	UpdateBody string
	DeleteBody string

	// PandoraIR, when non-nil, provides the fully-typed schema (including nested
	// block expand/flatten models) resolved from the Pandora Data API. It is
	// populated by ResolveBlocks() when -arm-type (or -service + -resource) is
	// supplied, and drives the block model transform pass.
	PandoraIR *ir.ResourceIR

	// ExpandFlattenSrc is the rendered source of the expand/flatten helpers
	// generated from the Pandora IR blocks. It is appended to the generated file.
	ExpandFlattenSrc string

	// Imports from the original file (path -> alias, "" means no alias).
	Imports map[string]string

	// Internal AST handles.
	fset *token.FileSet
	file *ast.File
	src  []byte
}

// Analyze parses the file at path and builds a full Info for typed-upgrade.
// The file must contain an untyped (native Plugin SDK) resource.
func Analyze(path string) (*Info, error) {
	base, err := list_upgrade.Analyze(path)
	if err != nil {
		return nil, err
	}
	if base.Kind != list_upgrade.KindUntyped {
		return nil, fmt.Errorf("file %q is not an untyped resource (kind: %s); typed-upgrade only applies to native Plugin SDK resources", path, base.Kind)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	info := &Info{
		Path:           path,
		Package:        base.Package,
		BaseName:       base.BaseName,
		GoName:         base.BaseName,
		StructName:     base.BaseName + "Resource",
		ModelName:      base.BaseName + "ResourceModel",
		TerraformType:  deriveTerraformType(path, base.ConstructorFunc, base.BaseName),
		SDKPackage:     base.SDKPackage,
		SDKImportPath:  base.SDKImportPath,
		IDPackage:      base.IDPackage,
		IDImportPath:   base.IDImportPath,
		IDTypeName:     base.IDTypeName,
		IDBase:         base.IDBase,
		IDParseFunc:    base.IDParseFunc,
		IDValidateFunc: base.IDValidateFunc,
		ServiceName:    base.ServiceName,
		ClientField:    base.ClientField,
		CreateFunc:     base.CreateFunc,
		ReadFunc:       base.ReadFunc,
		UpdateFunc:     base.UpdateFunc,
		DeleteFunc:     base.DeleteFunc,
		HasUpdate:      base.UpdateFunc != "",
		Imports:        map[string]string{},
		fset:           fset,
		file:           file,
		src:            src,
	}

	info.collectImports()

	funcs := collectFuncs(file)
	ctor := findConstructor(file)
	if ctor == nil {
		return nil, fmt.Errorf("could not find the pluginsdk.Resource constructor in %q", path)
	}

	resourceLit := returnedResourceLit(ctor)
	if resourceLit != nil {
		info.Fields = extractSchemaFields(src, fset, resourceLit)
		info.extractTimeouts(resourceLit)
		info.extractFlags(resourceLit)
	}

	info.extractCRUDBodies(funcs)
	info.applyModelTransforms()

	return info, nil
}

// collectImports indexes the file's imports by path → alias.
func (info *Info) collectImports() {
	for _, imp := range info.file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		info.Imports[path] = alias
	}
}

// extractTimeouts reads the Timeouts field of the resource literal.
func (info *Info) extractTimeouts(resourceLit *ast.CompositeLit) {
	for _, el := range resourceLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Timeouts" {
			continue
		}
		lit := unwrapSchemaLit(kv.Value)
		if lit == nil {
			break
		}
		for _, el2 := range lit.Elts {
			kv2, ok := el2.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			opKey, ok := kv2.Key.(*ast.Ident)
			if !ok {
				continue
			}
			dur := info.extractTimeoutDuration(kv2.Value)
			switch opKey.Name {
			case "Create":
				info.CreateTimeout = dur
			case "Read":
				info.ReadTimeout = dur
			case "Update":
				info.UpdateTimeout = dur
			case "Delete":
				info.DeleteTimeout = dur
			}
		}
		break
	}
	// Defaults when not found.
	if info.CreateTimeout == "" {
		info.CreateTimeout = "30 * time.Minute"
	}
	if info.ReadTimeout == "" {
		info.ReadTimeout = "5 * time.Minute"
	}
	if info.UpdateTimeout == "" && info.HasUpdate {
		info.UpdateTimeout = "30 * time.Minute"
	}
	if info.DeleteTimeout == "" {
		info.DeleteTimeout = "30 * time.Minute"
	}
}

// extractTimeoutDuration extracts the duration argument from
// pluginsdk.DefaultTimeout(X) as raw source text.
func (info *Info) extractTimeoutDuration(expr ast.Expr) string {
	call, ok := expr.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return "30 * time.Minute"
	}
	return nodeText(info.src, info.fset, call.Args[0].Pos(), call.Args[0].End())
}

// extractFlags reads CustomizeDiff, SchemaVersion, and DeprecationMessage
// from the resource literal.
func (info *Info) extractFlags(resourceLit *ast.CompositeLit) {
	for _, el := range resourceLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "CustomizeDiff":
			info.HasCustomizeDiff = true
		case "SchemaVersion":
			info.HasSchemaVersion = true
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.INT {
				fmt.Sscanf(lit.Value, "%d", &info.SchemaVersion)
			}
		case "DeprecationMessage":
			if lit, ok := kv.Value.(*ast.BasicLit); ok {
				info.DeprecationMsg = strings.Trim(lit.Value, `"`)
			}
		}
	}
}

// extractCRUDBodies extracts and transforms the body of each CRUD function.
func (info *Info) extractCRUDBodies(funcs map[string]*ast.FuncDecl) {
	if fn := funcs[info.CreateFunc]; fn != nil {
		info.CreateBody = info.transformBody(fn)
	}
	if fn := funcs[info.ReadFunc]; fn != nil {
		info.ReadBody = info.transformBody(fn)
	}
	if info.HasUpdate {
		if fn := funcs[info.UpdateFunc]; fn != nil {
			info.UpdateBody = info.transformBody(fn)
		}
	}
	if fn := funcs[info.DeleteFunc]; fn != nil {
		info.DeleteBody = info.transformBody(fn)
	}
}

// transformBody extracts the body of a CRUD function and applies safe
// text-level migrations to make it compatible with the typed SDK wrapper.
func (info *Info) transformBody(fn *ast.FuncDecl) string {
	if fn.Body == nil {
		return ""
	}
	// Extract between the opening and closing braces (exclusive).
	start := info.fset.Position(fn.Body.Lbrace).Offset + 1
	end := info.fset.Position(fn.Body.Rbrace).Offset
	if start >= end || end > len(info.src) {
		return ""
	}
	body := string(info.src[start:end])
	body = applyMigrations(body, info.ReadFunc)
	return body
}

// applyMigrations applies safe text-level transformations to a CRUD function
// body, converting native pluginsdk patterns to the typed SDK equivalents.
func applyMigrations(body, readFuncName string) string {
	lines := strings.Split(body, "\n")
	var out []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove timeout context setup lines.
		if strings.Contains(trimmed, "timeouts.For") && strings.Contains(trimmed, "cancel") {
			continue
		}
		if trimmed == "defer cancel()" {
			continue
		}

		// Remove the return-read call (typed wrapper auto-reads after Create/Update).
		if readFuncName != "" && strings.Contains(trimmed, "return "+readFuncName+"(") {
			line = replaceInLine(line, true, "return nil")
		}

		// Handle `d.SetId("")` followed by `return nil` on the next line:
		// replace the two-line pattern with a single `return metadata.MarkAsGone(id)`.
		if strings.Contains(trimmed, `d.SetId("")`) {
			indent := line[:len(line)-len(strings.TrimLeft(line, "\t "))]
			// Look ahead: if the next non-empty line is `return nil`, consume it.
			nextIsReturnNil := false
			for j := i + 1; j < len(lines); j++ {
				nt := strings.TrimSpace(lines[j])
				if nt == "" {
					continue
				}
				if nt == "return nil" {
					nextIsReturnNil = true
					lines[j] = "" // suppress the next return nil
				}
				break
			}
			if nextIsReturnNil {
				out = append(out, indent+"return metadata.MarkAsGone(id) // TODO: verify id variable name")
			} else {
				out = append(out, indent+"return metadata.MarkAsGone(id) // TODO: verify id variable name")
			}
			continue
		}
		// Suppress a return nil that was already consumed by the MarkAsGone look-ahead.
		if trimmed == "" && line == "" {
			out = append(out, line)
			continue
		}

		// Replace client access.
		line = strings.ReplaceAll(line, "meta.(*clients.Client).", "metadata.Client.")
		// Replace ResourceData accessors.
		line = replaceAll(line,
			"d.SetId(", "metadata.ResourceData.SetId(",
			"d.Id()", "metadata.ResourceData.Id()",
			"d.Set(", "metadata.ResourceData.Set(",
			"d.Get(", "metadata.ResourceData.Get(",
			"d.GetOk(", "metadata.ResourceData.GetOk(",
			"d.GetRawState(", "metadata.ResourceData.GetRawState(",
			"d.GetRawConfig(", "metadata.ResourceData.GetRawConfig(",
			"d.GetRawPlan(", "metadata.ResourceData.GetRawPlan(",
			"d.HasChange(", "metadata.ResourceData.HasChange(",
			"d.HasChanges(", "metadata.ResourceData.HasChanges(",
			"d.IsNewResource()", "metadata.ResourceData.IsNewResource()",
			"d.Timeout(", "metadata.ResourceData.Timeout(",
			// bare `d` as a function argument
			"(d, ", "(metadata.ResourceData, ",
			"(d)", "(metadata.ResourceData)",
			", d)", ", metadata.ResourceData)",
			", d, ", ", metadata.ResourceData, ",
		)
		// Replace bare `meta` used as a *clients.Client argument
		// (e.g. sdk.SetIDCallback(meta, ...) or subscriptionId := meta.(*clients.Client).Account.SubscriptionId).
		// The meta.(*clients.Client). prefix was already handled above; here we catch
		// plain `meta` remaining as a function argument.
		line = replaceAll(line,
			"(meta, ", "(metadata.Client, ",
			", meta, ", ", metadata.Client, ",
			", meta)", ", metadata.Client)",
		)

		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// replaceInLine conditionally replaces a whole line's read-call with return nil.
func replaceInLine(line string, cond bool, replacement string) string {
	if !cond {
		return line
	}
	// Preserve indentation.
	indent := line[:len(line)-len(strings.TrimLeft(line, "\t "))]
	return indent + replacement
}

// replaceAll applies successive pair-wise replacements to s.
func replaceAll(s string, pairs ...string) string {
	for i := 0; i+1 < len(pairs); i += 2 {
		s = strings.ReplaceAll(s, pairs[i], pairs[i+1])
	}
	return s
}

// collectFuncs indexes all top-level function declarations by name.
func collectFuncs(file *ast.File) map[string]*ast.FuncDecl {
	funcs := make(map[string]*ast.FuncDecl)
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv == nil {
			funcs[fd.Name.Name] = fd
		}
	}
	return funcs
}

// findConstructor returns the constructor function declaration for a native
// pluginsdk.Resource (a top-level func returning *pluginsdk.Resource).
func findConstructor(file *ast.File) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil {
			continue
		}
		if isPluginSDKResourceReturn(fd) {
			return fd
		}
	}
	return nil
}

// isPluginSDKResourceReturn reports whether fn returns *pluginsdk.Resource.
func isPluginSDKResourceReturn(fn *ast.FuncDecl) bool {
	if fn.Type == nil || fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}
	result := fn.Type.Results.List[0]
	star, ok := result.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "pluginsdk" && sel.Sel.Name == "Resource"
}

// returnedResourceLit returns the *pluginsdk.Resource composite literal
// returned by the constructor function.
func returnedResourceLit(fn *ast.FuncDecl) *ast.CompositeLit {
	var out *ast.CompositeLit
	ast.Inspect(fn, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		expr := ret.Results[0]
		if u, ok := expr.(*ast.UnaryExpr); ok {
			expr = u.X
		}
		if cl, ok := expr.(*ast.CompositeLit); ok {
			out = cl
			return false
		}
		return true
	})
	return out
}

// deriveTerraformType resolves the Terraform resource type string. It first
// looks in the sibling registration.go for a SupportedResources() entry whose
// constructor call matches constructorName (e.g. "resourceKustoDatabaseScript"
// → "azurerm_kusto_script"). Falls back to the snake-cased base name when the
// registration file cannot be found or parsed.
func deriveTerraformType(resourcePath, constructorName, baseName string) string {
	regPath := deriveRegistrationPath(resourcePath)
	if regPath != "" {
		if tfType := registrationKeyForConstructor(regPath, constructorName); tfType != "" {
			return tfType
		}
	}
	return "azurerm_" + toSnake(baseName)
}

// registrationKeyForConstructor parses the SupportedResources() map in
// registration.go and returns the map key whose value calls constructorName.
func registrationKeyForConstructor(regPath, constructorName string) string {
	src, err := os.ReadFile(regPath)
	if err != nil {
		return ""
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, regPath, src, 0)
	if err != nil {
		return ""
	}
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Name.Name != "SupportedResources" {
			continue
		}
		var found string
		ast.Inspect(fd, func(n ast.Node) bool {
			if found != "" {
				return false
			}
			mapLit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			for _, el := range mapLit.Elts {
				kv, ok := el.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.BasicLit)
				if !ok || key.Kind != token.STRING {
					continue
				}
				if callFuncName(callExprFromValue(kv.Value)) == constructorName {
					found = strings.Trim(key.Value, `"`)
					return false
				}
			}
			return true
		})
		if found != "" {
			return found
		}
	}
	return ""
}

// callExprFromValue unwraps a call expression (possibly behind a unary & or
// address-of) from a map value. It returns nil when the value is not a call.
func callExprFromValue(expr ast.Expr) *ast.CallExpr {
	if call, ok := expr.(*ast.CallExpr); ok {
		return call
	}
	return nil
}

// toSnake converts a PascalCase or camelCase identifier to snake_case.
func toSnake(s string) string {
	var out []rune
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 && (runes[i-1] < 'A' || runes[i-1] > 'Z') {
				out = append(out, '_')
			} else if i > 0 && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				out = append(out, '_')
			}
			out = append(out, r+'a'-'A')
		} else {
			out = append(out, r)
		}
	}
	result := string(out)
	// Collapse multiple underscores.
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.TrimPrefix(result, "_")
}
