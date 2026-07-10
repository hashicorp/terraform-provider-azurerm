package list_upgrade

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// addIdentityAssertionEdit ensures a `var _ sdk.ResourceWithIdentity = X{}`
// interface assertion is present, either by extending an existing `var (...)`
// block or by promoting a single assertion to a block.
func (r *Resource) addIdentityAssertionEdit(e *editor) {
	gd := r.assertDecl
	line := fmt.Sprintf("_ sdk.ResourceWithIdentity = %s{}", r.StructName)

	if gd.Lparen.IsValid() {
		// Existing block: insert a new spec right after the opening paren.
		e.insert(r.offset(gd.Lparen)+1, "\n\t"+line)
		return
	}

	// Single assertion: rewrite the whole declaration as a block that keeps the
	// existing spec and adds the identity assertion first.
	existing := r.nodeText(gd.Specs[0].Pos(), gd.Specs[0].End())
	block := fmt.Sprintf("var (\n\t%s\n\t%s\n)", line, existing)
	e.replace(r.offset(gd.Pos()), r.offset(gd.End()), block)
}

// addIdentityMethodEdit inserts the Identity() method returning the resource's
// ID type, immediately after the interface-assertion declaration.
func (r *Resource) addIdentityMethodEdit(e *editor) {
	method := fmt.Sprintf("\n\nfunc (r %s) Identity() resourceids.ResourceId {\n\treturn &%s.%s{}\n}\n",
		r.StructName, r.SDKPackage, r.IDTypeName)
	e.insert(r.offset(r.assertDecl.End()), method)
}

// addCreateIdentityEdit injects a SetResourceIdentityData call into Create,
// right after metadata.SetID(id). When Create ends with `return nil` after
// SetID that return is replaced; otherwise an error-checked call is inserted.
func (r *Resource) addCreateIdentityEdit(e *editor) error {
	create := r.methods["Create"]
	if create == nil {
		return fmt.Errorf("resource has no Create method")
	}
	body := resourceFuncBody(create)
	if body == nil || body.Body == nil {
		return fmt.Errorf("could not locate the Create closure body")
	}

	stmts := body.Body.List
	for i, st := range stmts {
		idVar, ok := setIDCallVar(st)
		if !ok {
			continue
		}
		if i+1 < len(stmts) && isReturnNil(stmts[i+1]) {
			e.replace(r.offset(stmts[i+1].Pos()), r.offset(stmts[i+1].End()),
				fmt.Sprintf("return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &%s)", idVar))
			return nil
		}
		insert := fmt.Sprintf("\nif err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &%s); err != nil {\nreturn err\n}\n", idVar)
		e.insert(r.offset(st.End()), insert)
		return nil
	}
	return fmt.Errorf("could not find a metadata.SetID(...) call in Create")
}

// setIDCallVar reports whether st is `metadata.SetID(<var>)` and returns the
// argument variable name.
func setIDCallVar(st ast.Stmt) (string, bool) {
	es, ok := st.(*ast.ExprStmt)
	if !ok {
		return "", false
	}
	call, ok := es.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return "", false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "SetID" {
		return "", false
	}
	if x, ok := sel.X.(*ast.Ident); !ok || x.Name != "metadata" {
		return "", false
	}
	if id, ok := call.Args[0].(*ast.Ident); ok {
		return id.Name, true
	}
	return "", false
}

// addGoGenerateEdit inserts the resource-identity test go:generate directive
// beneath the import block, unless one is already present.
func (r *Resource) addGoGenerateEdit(e *editor, resourceName, properties string) {
	if r.hasResourceIdentityGoGenerate() {
		return
	}
	importDecl := r.importDecl()
	if importDecl == nil {
		return
	}
	line := fmt.Sprintf("\n\n//go:generate go run ../../tools/generator-tests resourceidentity -resource-name %s -properties %q",
		resourceName, properties)
	e.insert(r.offset(importDecl.End()), line)
}

// hasResourceIdentityGoGenerate reports whether the file already declares a
// resourceidentity go:generate directive.
func (r *Resource) hasResourceIdentityGoGenerate() bool {
	for _, cg := range r.file.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "//go:generate") && strings.Contains(c.Text, "resourceidentity") {
				return true
			}
		}
	}
	return false
}

// importDecl returns the file's import declaration, if any.
func (r *Resource) importDecl() *ast.GenDecl {
	for _, decl := range r.file.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			return gd
		}
	}
	return nil
}

// injectIdentityIntoExistingFlattenEdit adds a SetResourceIdentityData call
// before the terminal encode of an existing flatten method.
func (r *Resource) injectIdentityIntoExistingFlattenEdit(e *editor) error {
	flatten := r.methods["flatten"]
	if flatten == nil || flatten.Body == nil {
		return fmt.Errorf("resource has no flatten method to update")
	}
	for _, st := range flatten.Body.List {
		if pos, _, ok := encodeStatementAny(st); ok {
			e.insert(r.offset(pos),
				"if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {\nreturn err\n}\n\n")
			return nil
		}
	}
	return fmt.Errorf("could not find a terminal metadata.Encode in the flatten method")
}

// encodeStatementAny is encodeStatement without pinning the state variable name,
// used for existing flatten methods where the variable name is unknown.
func encodeStatementAny(st ast.Stmt) (start, end token.Pos, ok bool) {
	switch s := st.(type) {
	case *ast.ReturnStmt:
		if len(s.Results) == 1 && isEncodeCallAny(s.Results[0]) {
			return s.Pos(), s.End(), true
		}
	case *ast.IfStmt:
		if as, ok := s.Init.(*ast.AssignStmt); ok && len(as.Rhs) == 1 && isEncodeCallAny(as.Rhs[0]) {
			return s.Pos(), s.End(), true
		}
	}
	return 0, 0, false
}

// isEncodeCallAny reports whether expr is `metadata.Encode(...)`.
func isEncodeCallAny(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Encode" {
		return false
	}
	x, ok := sel.X.(*ast.Ident)
	return ok && x.Name == "metadata"
}
