// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package list_upgrade

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// analyzeUntyped derives the metadata for a native Plugin SDK resource from its
// `func resourceX() *pluginsdk.Resource` constructor.
func (r *Resource) analyzeUntyped(ctor *ast.FuncDecl) {
	r.ConstructorFunc = ctor.Name.Name
	r.BaseName = strings.TrimPrefix(ctor.Name.Name, "resource")

	r.resourceLit = returnedResourceLiteral(ctor)
	if r.resourceLit != nil {
		for _, el := range r.resourceLit.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			switch key.Name {
			case "Create":
				r.CreateFunc = identName(kv.Value)
			case "Read":
				r.ReadFunc = identName(kv.Value)
			case "Update":
				r.UpdateFunc = identName(kv.Value)
			case "Delete":
				r.DeleteFunc = identName(kv.Value)
			}
		}
	}

	// The Read func is the source of truth for the ID, client and Get method.
	if read := r.funcs[r.ReadFunc]; read != nil {
		if pkg, base := parseIDCall(read); base != "" {
			r.IDBase = base
			r.IDPackage = pkg
			if imp, ok := r.imports[pkg]; ok {
				r.IDImportPath = strings.Trim(imp.Path.Value, `"`)
			}
		}
		if svc, field := untypedClientAccessor(read); field != "" {
			r.ServiceName = svc
			r.ClientField = field
		}
		r.GetMethod = untypedGetMethod(read)
	}
	if r.IDBase != "" {
		r.IDTypeName = r.IDBase + "Id"
		r.IDParseFunc = "Parse" + r.IDBase + "ID"
		r.IDValidateFunc = "Validate" + r.IDBase + "ID"
	}

	// The go-azure-sdk resource-manager import provides the model + client. When
	// several are imported (e.g. subnet imports subnets, serviceendpointpolicies
	// and ipampools) the client field / Get call disambiguates which one owns the
	// Get method, so this runs after the Read func has been parsed.
	if pkg, path := r.resourceManagerImport(); pkg != "" {
		r.SDKPackage = pkg
		r.SDKImportPath = path
	}
	// If the ID lives in the same package as the model, keep SDKPackage as the
	// ID package too so callers that only read SDKPackage still work.
	if r.SDKPackage == "" {
		r.SDKPackage = r.IDPackage
		r.SDKImportPath = r.IDImportPath
	}

	// The authoritative read model is resp.Model's type, read from the vendored
	// SDK Get method rather than inferred from Pandora.
	r.deriveReadModel()

	r.detectUntypedParent()
	r.deriveListMethods()

	r.FlattenFunc = r.ConstructorFunc + "Flatten"
	_, r.HasFlatten = r.funcs[r.FlattenFunc]
	if !r.HasFlatten {
		r.FlattenFunc = ""
	} else {
		r.FlattenIDValue = flattenIDIsValue(r.funcs[r.FlattenFunc])
	}

	r.HasIdentity = r.hasIdentityField() && r.hasIdentityImporter()
}

// TestStructName reads the sibling <base>_resource_test.go file and returns the
// name of the acceptance-test struct (the one exposing a basic(data) method), so
// generated list tests reference the real struct even when its casing differs
// from the resource base name (e.g. VirtualHubIpResource vs VirtualHubIP).
func TestStructName(resourcePath string) string {
	testPath := strings.TrimSuffix(resourcePath, "_resource.go") + "_resource_test.go"
	src, err := os.ReadFile(testPath)
	if err != nil {
		return ""
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, testPath, src, 0)
	if err != nil {
		return ""
	}
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 || fd.Name.Name != "basic" {
			continue
		}
		if name := receiverTypeName(fd.Recv.List[0].Type); name != "" {
			return name
		}
	}
	return ""
}

// returnedResourceLiteral returns the &pluginsdk.Resource{...} composite literal
// returned by the constructor.
func returnedResourceLiteral(ctor *ast.FuncDecl) *ast.CompositeLit {
	var out *ast.CompositeLit
	ast.Inspect(ctor, func(n ast.Node) bool {
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

// identName returns the identifier name of an expression, or "" when it is not a
// bare identifier.
func identName(expr ast.Expr) string {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// resourceManagerImport returns the package name and path of the go-azure-sdk
// resource-manager import that provides the model + client. When a file imports
// several such packages (e.g. subnet imports subnets, serviceendpointpolicies
// and ipampools), it selects the one that owns the Get method — preferring the
// package referenced in the Get call's arguments, then the one matching the
// client field name — and falls back to the first alphabetically so selection is
// deterministic.
func (r *Resource) resourceManagerImport() (pkg, path string) {
	type rmImport struct{ name, path string }
	var imports []rmImport
	names := map[string]bool{}
	for name, imp := range r.imports {
		p := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(p, "/go-azure-sdk/resource-manager/") {
			imports = append(imports, rmImport{name, p})
			names[name] = true
		}
	}
	if len(imports) == 0 {
		return "", ""
	}
	sort.Slice(imports, func(i, j int) bool { return imports[i].name < imports[j].name })
	if len(imports) == 1 {
		return imports[0].name, imports[0].path
	}

	pick := getCallPackageQualifier(r.funcs[r.ReadFunc], names)
	if pick == "" && r.ClientField != "" {
		want := strings.ToLower(r.ClientField)
		for _, im := range imports {
			if strings.ToLower(im.name) == want {
				pick = im.name
				break
			}
		}
	}
	for _, im := range imports {
		if im.name == pick {
			return im.name, im.path
		}
	}
	return imports[0].name, imports[0].path
}

// getCallPackageQualifier returns the package qualifier used in an argument of
// the Read func's Get call (e.g. `subnets.DefaultGetOperationOptions()` yields
// "subnets"), restricted to one of the candidate package names. This
// disambiguates the SDK package when several are imported.
func getCallPackageQualifier(read *ast.FuncDecl, candidates map[string]bool) string {
	if read == nil || read.Body == nil {
		return ""
	}
	respVar, _, idx := findGetAssignment(read.Body.List, read)
	if respVar == "" || idx < 0 {
		return ""
	}
	as, ok := read.Body.List[idx].(*ast.AssignStmt)
	if !ok || len(as.Rhs) != 1 {
		return ""
	}
	call, ok := as.Rhs[0].(*ast.CallExpr)
	if !ok {
		return ""
	}
	var found string
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if id, ok := sel.X.(*ast.Ident); ok && candidates[id.Name] {
				found = id.Name
				return false
			}
			return true
		})
		if found != "" {
			break
		}
	}
	return found
}

// flattenIDIsValue reports whether the flatten function takes its id parameter by
// value (e.g. `id commonids.SubnetId`) rather than by pointer (`id *X`). The id
// is the second parameter, after the ResourceData.
func flattenIDIsValue(fn *ast.FuncDecl) bool {
	if fn == nil || fn.Type == nil || fn.Type.Params == nil {
		return false
	}
	idx := 0
	for _, p := range fn.Type.Params.List {
		n := len(p.Names)
		if n == 0 {
			n = 1
		}
		for k := 0; k < n; k++ {
			if idx == 1 { // second parameter is the id
				_, isPtr := p.Type.(*ast.StarExpr)
				return !isPtr
			}
			idx++
		}
	}
	return false
}

// untypedClientAccessor extracts service and field from a
// `meta.(*clients.Client).<Service>.<Field>` selector in the given function.
func untypedClientAccessor(fd *ast.FuncDecl) (service, field string) {
	ast.Inspect(fd, func(n ast.Node) bool {
		outer, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		mid, ok := outer.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ta, ok := mid.X.(*ast.TypeAssertExpr)
		if !ok {
			return true
		}
		if id, ok := ta.X.(*ast.Ident); !ok || id.Name != "meta" {
			return true
		}
		service = mid.Sel.Name
		field = outer.Sel.Name
		return false
	})
	return service, field
}

// untypedGetMethod returns the client method whose result is later read via
// `<var>.Model` (i.e. the resource's Get call), e.g. "Get" or
// "VirtualHubIPConfigurationGet".
func untypedGetMethod(read *ast.FuncDecl) string {
	if read.Body == nil {
		return ""
	}
	_, method, _ := findGetAssignment(read.Body.List, read)
	return method
}

// findGetAssignment locates the `resp, err := client.<Method>(...)` statement
// whose result is subsequently accessed as `resp.Model`, returning the response
// variable name, the method name and the statement index. stmts is the list of
// statements to scan and scope is the node to inspect for `.Model` accesses
// (the function or closure containing them).
func findGetAssignment(stmts []ast.Stmt, scope ast.Node) (respVar, method string, index int) {
	// Collect variables that are later accessed via `.Model`.
	modelVars := map[string]bool{}
	ast.Inspect(scope, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel.Name == "Model" {
			if id, ok := sel.X.(*ast.Ident); ok {
				modelVars[id.Name] = true
			}
		}
		return true
	})
	for i, st := range stmts {
		as, ok := st.(*ast.AssignStmt)
		if !ok || len(as.Lhs) != 2 || len(as.Rhs) != 1 {
			continue
		}
		lhs, ok := as.Lhs[0].(*ast.Ident)
		if !ok || !modelVars[lhs.Name] {
			continue
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok {
			continue
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			return lhs.Name, sel.Sel.Name, i
		}
	}
	return "", "", -1
}

// hasIdentityField reports whether the resource literal declares an Identity field.
func (r *Resource) hasIdentityField() bool {
	if r.resourceLit == nil {
		return false
	}
	for _, el := range r.resourceLit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "Identity" {
				return true
			}
		}
	}
	return false
}

// hasIdentityImporter reports whether the Importer uses ImporterValidatingIdentity.
func (r *Resource) hasIdentityImporter() bool {
	if r.resourceLit == nil {
		return false
	}
	for _, el := range r.resourceLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if key, ok := kv.Key.(*ast.Ident); !ok || key.Name != "Importer" {
			continue
		}
		if call, ok := kv.Value.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "ImporterValidatingIdentity" {
				return true
			}
		}
	}
	return false
}

// upgradeUntyped computes the edits for a native Plugin SDK resource.
func (r *Resource) upgradeUntyped(opts UpgradeOptions) (newSrc []byte, changed bool, err error) {
	e := newEditor(r.src)
	withIdentity := r.HasIdentity || opts.AddIdentity

	if opts.ExtractFlatten && !r.HasFlatten {
		if err := r.extractUntypedFlattenEdits(e, opts.ReadModel, withIdentity); err != nil {
			return nil, false, fmt.Errorf("extracting flatten: %w", err)
		}
	}

	if opts.AddIdentity && !r.HasIdentity {
		if err := r.addUntypedIdentityEdits(e, opts); err != nil {
			return nil, false, fmt.Errorf("adding identity: %w", err)
		}
	}

	if !e.hasEdits() {
		return r.src, false, nil
	}
	out, err := e.bytes()
	if err != nil {
		return nil, false, err
	}
	return out, true, nil
}

// extractUntypedFlattenEdits moves the state-setting region of Read into a
// dedicated resourceXFlatten function and rewires Read to call it.
func (r *Resource) extractUntypedFlattenEdits(e *editor, readModel string, withIdentity bool) error {
	if r.SDKPackage == "" || r.IDTypeName == "" {
		return fmt.Errorf("missing SDK package or ID type; cannot build the flatten signature")
	}
	if readModel == "" {
		return fmt.Errorf("the SDK read model type is required to build the flatten signature")
	}
	read := r.funcs[r.ReadFunc]
	if read == nil || read.Body == nil {
		return fmt.Errorf("could not locate the Read function %q", r.ReadFunc)
	}

	respVar, _, idx := findGetAssignment(read.Body.List, read)
	if idx == -1 {
		return fmt.Errorf("could not find the Get call in %q", r.ReadFunc)
	}
	stmts := read.Body.List
	// The statement after the Get call should be its error handler; the region
	// to extract is everything after that to the end of the function.
	regionStart := idx + 2
	if regionStart >= len(stmts) {
		return fmt.Errorf("no state-setting statements to extract in %q", r.ReadFunc)
	}
	startPos := stmts[regionStart].Pos()
	endPos := stmts[len(stmts)-1].End()

	modelVar := modelVarInRegion(stmts[regionStart:])
	if modelVar == "" {
		modelVar = "model"
	}

	regionText := string(r.src[r.offset(startPos):r.offset(endPos)])
	body := rewriteRespModelVar(regionText, respVar, modelVar)
	if strings.Contains(body, respVar+".") {
		return fmt.Errorf("the flatten region references %q beyond `.Model`; refactor Read manually (see guide-list-resource.md step 1)", respVar)
	}

	idPkg := r.IDPackage
	if idPkg == "" {
		idPkg = r.SDKPackage
	}

	flattenName := r.ConstructorFunc + "Flatten"
	var fn strings.Builder
	fmt.Fprintf(&fn, "\n\nfunc %s(d *pluginsdk.ResourceData, id *%s.%s, %s *%s.%s) error {\n",
		flattenName, idPkg, r.IDTypeName, modelVar, r.SDKPackage, readModel)
	if withIdentity && !strings.Contains(body, "SetResourceIdentityData") {
		body = injectIdentityBeforeFinalReturn(body)
	}
	fn.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		fn.WriteString("\n")
	}
	fn.WriteString("}\n")

	// Insert the flatten func immediately after the Read function.
	e.insert(r.offset(read.End()), fn.String())
	// Replace the extracted region with a delegating call.
	e.replace(r.offset(startPos), r.offset(endPos), fmt.Sprintf("return %s(d, id, %s.Model)", flattenName, respVar))
	return nil
}

// rewriteRespModelVar rewrites `if <mv> := <resp>.Model; <mv> != nil {` to
// `if <mv> != nil {` and any bare `<resp>.Model` to the model parameter name.
func rewriteRespModelVar(text, respVar, modelVar string) string {
	text = strings.ReplaceAll(text, fmt.Sprintf("if %s := %s.Model; ", modelVar, respVar), "if ")
	text = strings.ReplaceAll(text, respVar+".Model", modelVar)
	return text
}

// injectIdentityBeforeFinalReturn inserts a SetResourceIdentityData call before
// the final return statement of an untyped flatten body.
func injectIdentityBeforeFinalReturn(body string) string {
	idx := strings.LastIndex(body, "return ")
	if idx == -1 {
		return body + "\nif err := pluginsdk.SetResourceIdentityData(d, id); err != nil {\nreturn err\n}\n"
	}
	inject := "if err := pluginsdk.SetResourceIdentityData(d, id); err != nil {\nreturn err\n}\n\n"
	return body[:idx] + inject + body[idx:]
}
