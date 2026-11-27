// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// generator-client-providers generates lightweight per-service client Provider interfaces
// and From(meta any) accessors in internal/services/*/client/provider_gen.go (plus
// internal/preflight/client and internal/clients/graph), and their corresponding getter
// methods and compile-time interface checks in internal/clients/provider_gen.go.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const modulePrefix = "github.com/hashicorp/terraform-provider-azurerm/"

type FieldMapping struct {
	FieldName        string `json:"field_name"`
	GetterName       string `json:"getter_name"`
	InterfaceName    string `json:"interface_name"`
	HelperFuncName   string `json:"helper_func_name"`
	TargetImportPath string `json:"target_import_path"`
	TargetPkgName    string `json:"target_pkg_name"`
	ClientsTypeExpr  string `json:"clients_type_expr"`
	ClientsTypePkg   string `json:"clients_type_pkg"`
	LocalTypeExpr    string `json:"local_type_expr"`
	ExternalTypeImp  string `json:"external_type_imp,omitempty"`
	ExternalTypeName string `json:"external_type_name,omitempty"`
}

func main() {
	repoRootFlag := flag.String("repo-root", "", "path to repository root (auto-detected if empty)")
	dumpJSONFlag := flag.String("dump-json", "", "optional path to write field mappings JSON")

	flag.Parse()

	repoRoot := *repoRootFlag
	if repoRoot == "" {
		var err error

		repoRoot, err = detectRepoRoot()
		if err != nil {
			fatalf("detecting repo root: %v", err)
		}
	}

	mappings, clientsImps, err := extractMappings(repoRoot)
	if err != nil {
		fatalf("extracting client mappings: %v", err)
	}

	if *dumpJSONFlag != "" {
		data, err := json.MarshalIndent(mappings, "", "  ")
		if err != nil {
			fatalf("marshaling json: %v", err)
		}

		if err := os.WriteFile(*dumpJSONFlag, data, 0o644); err != nil {
			fatalf("writing json: %v", err)
		}
	}

	if err := writeServiceProviders(repoRoot, mappings); err != nil {
		fatalf("writing service provider_gen.go files: %v", err)
	}

	if err := writeClientsProviderGen(repoRoot, mappings, clientsImps); err != nil {
		fatalf("writing internal/clients/provider_gen.go: %v", err)
	}
}

func fatalf(formatStr string, args ...any) {
	fmt.Fprintf(os.Stderr, "generator-client-providers: "+formatStr+"\n", args...)
	os.Exit(1)
}

func detectRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if filepath.Base(dir) != ".tools" {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}

		dir = parent
	}
}

func extractMappings(repoRoot string) ([]FieldMapping, map[string]string, error) {
	clientsDir := filepath.Join(repoRoot, "internal", "clients")
	fset := token.NewFileSet()

	allImports := make(map[string]string) // alias -> import path

	type rawField struct {
		fieldName string
		isPtr     bool
		pkgAlias  string
		typeName  string
	}

	var fields []rawField

	buildPkgByField := make(map[string]string) // fieldName -> pkgAlias called in Build / buildAutoClients

	for _, fname := range []string{"client.go", "client_gen.go"} {
		f, err := parser.ParseFile(fset, filepath.Join(clientsDir, fname), nil, 0)
		if err != nil {
			return nil, nil, err
		}

		fileImports := make(map[string]string)

		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)

			alias := filepath.Base(p)
			if imp.Name != nil {
				alias = imp.Name.Name
			}

			fileImports[alias] = p
			allImports[alias] = p
		}

		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}

				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok || (ts.Name.Name != "Client" && ts.Name.Name != "autoClient") {
						continue
					}

					st, ok := ts.Type.(*ast.StructType)
					if !ok {
						continue
					}

					for _, fld := range st.Fields.List {
						if len(fld.Names) == 0 {
							continue
						}

						name := fld.Names[0].Name
						if name == "StopContext" || name == "Account" || name == "Features" {
							continue
						}

						tExpr := fld.Type

						isPtr := false
						if star, ok := tExpr.(*ast.StarExpr); ok {
							isPtr = true
							tExpr = star.X
						}

						sel, ok := tExpr.(*ast.SelectorExpr)
						if !ok {
							continue
						}

						ident, ok := sel.X.(*ast.Ident)
						if !ok {
							continue
						}

						fields = append(fields, rawField{
							fieldName: name,
							isPtr:     isPtr,
							pkgAlias:  ident.Name,
							typeName:  sel.Sel.Name,
						})
					}
				}
			case *ast.FuncDecl:
				if d.Name.Name != "Build" && d.Name.Name != "buildAutoClients" {
					continue
				}

				ast.Inspect(d.Body, func(n ast.Node) bool {
					assign, ok := n.(*ast.AssignStmt)
					if !ok || len(assign.Lhs) == 0 || len(assign.Rhs) == 0 {
						return true
					}

					lhsSel, ok := assign.Lhs[0].(*ast.SelectorExpr)
					if !ok {
						return true
					}

					call, ok := assign.Rhs[0].(*ast.CallExpr)
					if !ok {
						return true
					}

					funSel, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}

					pkgIdent, ok := funSel.X.(*ast.Ident)
					if !ok {
						return true
					}

					buildPkgByField[lhsSel.Sel.Name] = fileImports[pkgIdent.Name]

					return true
				})
			}
		}
	}

	// Group by target package to detect multi-field client packages (like containers/client)
	countByTarget := make(map[string]int)
	targetByField := make(map[string]string)

	for _, rf := range fields {
		typeImp := allImports[rf.pkgAlias]

		targetImp := typeImp
		if !strings.HasPrefix(typeImp, modulePrefix) {
			if bImp, ok := buildPkgByField[rf.fieldName]; ok && strings.HasPrefix(bImp, modulePrefix) {
				targetImp = bImp
			} else {
				return nil, nil, fmt.Errorf("cannot determine local client package for field %s (type pkg %s)", rf.fieldName, typeImp)
			}
		}

		targetByField[rf.fieldName] = targetImp
		countByTarget[targetImp]++
	}

	var mappings []FieldMapping

	for _, rf := range fields {
		typeImp := allImports[rf.pkgAlias]
		targetImp := targetByField[rf.fieldName]
		relDir := strings.TrimPrefix(targetImp, modulePrefix)

		pkgName, err := readPackageName(filepath.Join(repoRoot, filepath.FromSlash(relDir)))
		if err != nil {
			return nil, nil, fmt.Errorf("reading package name for %s: %w", relDir, err)
		}

		ptrPrefix := ""
		if rf.isPtr {
			ptrPrefix = "*"
		}

		ifaceName := "Provider"
		helperName := "From"

		if countByTarget[targetImp] > 1 {
			if rf.typeName == "AutoClient" {
				ifaceName = "AutoProvider"
				helperName = "AutoFrom"
			} else if rf.typeName != "Client" {
				ifaceName = rf.fieldName + "Provider"
				helperName = rf.fieldName + "From"
			}
		}

		m := FieldMapping{
			FieldName:        rf.fieldName,
			GetterName:       rf.fieldName + "Client",
			InterfaceName:    ifaceName,
			HelperFuncName:   helperName,
			TargetImportPath: targetImp,
			TargetPkgName:    pkgName,
			ClientsTypePkg:   typeImp,
			ClientsTypeExpr:  ptrPrefix + rf.pkgAlias + "." + rf.typeName,
		}

		if typeImp == targetImp {
			m.LocalTypeExpr = ptrPrefix + rf.typeName
		} else {
			extAlias := rf.pkgAlias
			m.ExternalTypeImp = typeImp
			m.ExternalTypeName = extAlias
			m.LocalTypeExpr = ptrPrefix + extAlias + "." + rf.typeName
		}

		mappings = append(mappings, m)
	}

	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].FieldName < mappings[j].FieldName
	})

	return mappings, allImports, nil
}

func readPackageName(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}

		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, parser.PackageClauseOnly)
		if err == nil && f.Name != nil {
			return f.Name.Name, nil
		}
	}

	return filepath.Base(dir), nil
}

func writeServiceProviders(repoRoot string, mappings []FieldMapping) error {
	byTarget := make(map[string][]FieldMapping)
	for _, m := range mappings {
		byTarget[m.TargetImportPath] = append(byTarget[m.TargetImportPath], m)
	}

	for targetImp, list := range byTarget {
		sort.Slice(list, func(i, j int) bool {
			return list[i].InterfaceName < list[j].InterfaceName
		})

		relDir := strings.TrimPrefix(targetImp, modulePrefix)
		outPath := filepath.Join(repoRoot, filepath.FromSlash(relDir), "provider_gen.go")

		var buf bytes.Buffer
		buf.WriteString("// Copyright IBM Corp. 2014, 2026\n")
		buf.WriteString("// SPDX-License-Identifier: MPL-2.0\n\n")
		buf.WriteString("// Code generated by generator-client-providers; DO NOT EDIT.\n\n")
		fmt.Fprintf(&buf, "package %s\n\n", list[0].TargetPkgName)

		extImps := make(map[string]string)

		for _, m := range list {
			if m.ExternalTypeImp != "" {
				extImps[m.ExternalTypeName] = m.ExternalTypeImp
			}
		}

		if len(extImps) > 0 {
			var aliases []string
			for a := range extImps {
				aliases = append(aliases, a)
			}

			sort.Strings(aliases)
			buf.WriteString("import (\n")

			for _, a := range aliases {
				fmt.Fprintf(&buf, "\t%s %q\n", a, extImps[a])
			}

			buf.WriteString(")\n\n")
		}

		for _, m := range list {
			fmt.Fprintf(&buf, "type %s interface {\n\t%s() %s\n}\n\n", m.InterfaceName, m.GetterName, m.LocalTypeExpr)
			fmt.Fprintf(&buf, "func %s(meta any) %s {\n\treturn meta.(%s).%s()\n}\n\n", m.HelperFuncName, m.LocalTypeExpr, m.InterfaceName, m.GetterName)
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("formatting %s: %w\n%s", outPath, err, buf.String())
		}

		if err := os.WriteFile(outPath, formatted, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func writeClientsProviderGen(repoRoot string, mappings []FieldMapping, _ map[string]string) error {
	outPath := filepath.Join(repoRoot, "internal", "clients", "provider_gen.go")

	// Build deterministic unique aliases for all needed imports in internal/clients/provider_gen.go
	impToAlias := make(map[string]string)
	usedAliases := map[string]bool{
		"context":  true,
		"account":  true,
		"features": true,
	}

	registerImp := func(imp string) {
		if _, ok := impToAlias[imp]; ok {
			return
		}

		parts := strings.Split(imp, "/")

		base := parts[len(parts)-1]
		if base == "client" && len(parts) >= 2 {
			base = strings.ReplaceAll(parts[len(parts)-2], "-", "") + "Client"
		} else {
			base = strings.NewReplacer("-", "_", ".", "_").Replace(base)
			if len(parts) >= 2 && (strings.HasPrefix(base, "20") || strings.HasPrefix(base, "v20")) {
				base = strings.ReplaceAll(parts[len(parts)-2], "-", "") + "_" + base
			}
		}

		alias := base

		idx := 2
		for usedAliases[alias] {
			alias = fmt.Sprintf("%s%d", base, idx)
			idx++
		}

		usedAliases[alias] = true
		impToAlias[imp] = alias
	}

	for _, m := range mappings {
		registerImp(m.TargetImportPath)
		registerImp(m.ClientsTypePkg)
	}

	var impLines []string

	impLines = append(impLines, fmt.Sprintf("\t%q", modulePrefix+"internal/clients/account"))
	impLines = append(impLines, fmt.Sprintf("\t%q", modulePrefix+"internal/features"))

	for imp, alias := range impToAlias {
		if filepath.Base(imp) == alias {
			impLines = append(impLines, fmt.Sprintf("\t%q", imp))
		} else {
			impLines = append(impLines, fmt.Sprintf("\t%s %q", alias, imp))
		}
	}

	sort.Strings(impLines)

	var buf bytes.Buffer
	buf.WriteString("// Copyright IBM Corp. 2014, 2026\n")
	buf.WriteString("// SPDX-License-Identifier: MPL-2.0\n\n")
	buf.WriteString("// Code generated by generator-client-providers; DO NOT EDIT.\n\n")
	buf.WriteString("package clients\n\n")
	buf.WriteString("import (\n\t\"context\"\n\n")

	for _, l := range impLines {
		buf.WriteString(l + "\n")
	}

	buf.WriteString(")\n\n")

	buf.WriteString("var (\n")
	buf.WriteString("\t_ account.Provider  = (*Client)(nil)\n")
	buf.WriteString("\t_ features.Provider = (*Client)(nil)\n")

	for _, m := range mappings {
		tAlias := impToAlias[m.TargetImportPath]
		fmt.Fprintf(&buf, "\t_ %s.%s = (*Client)(nil)\n", tAlias, m.InterfaceName)
	}

	buf.WriteString(")\n\n")

	buf.WriteString("func (client *Client) ClientAccount() *account.ResourceManagerAccount {\n\treturn client.Account\n}\n\n")
	buf.WriteString("func (client *Client) ClientStopContext() context.Context {\n\treturn client.StopContext\n}\n\n")
	buf.WriteString("func (client *Client) UserFeatures() features.UserFeatures {\n\treturn client.Features\n}\n\n")

	for _, m := range mappings {
		typePkgAlias := impToAlias[m.ClientsTypePkg]
		ptr := ""

		rest := m.ClientsTypeExpr
		if strings.HasPrefix(rest, "*") {
			ptr = "*"
			rest = rest[1:]
		}

		dotIdx := strings.IndexByte(rest, '.')
		typeName := rest[dotIdx+1:]
		fmt.Fprintf(&buf, "func (client *Client) %s() %s%s.%s {\n\treturn client.%s\n}\n\n",
			m.GetterName, ptr, typePkgAlias, typeName, m.FieldName)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting %s: %w\n%s", outPath, err, buf.String())
	}

	return os.WriteFile(outPath, formatted, 0o644)
}
