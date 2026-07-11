// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package list_upgrade

import (
	"fmt"
	"go/ast"
	"strings"
)

// addUntypedIdentityEdits wires Resource Identity into a native Plugin SDK
// resource: it swaps the Importer for ImporterValidatingIdentity, adds the
// Identity schema field, sets the identity data in Create, and adds the
// identity-test go:generate directive.
func (r *Resource) addUntypedIdentityEdits(e *editor, opts UpgradeOptions) error {
	if r.resourceLit == nil {
		return fmt.Errorf("could not find the pluginsdk.Resource literal")
	}
	if r.IDTypeName == "" {
		return fmt.Errorf("could not determine the resource ID type")
	}
	idPkg := r.IDPackage
	if idPkg == "" {
		idPkg = r.SDKPackage
	}
	idRef := fmt.Sprintf("&%s.%s{}", idPkg, r.IDTypeName)

	imp := r.untypedImporterField()
	if imp == nil {
		return fmt.Errorf("could not find the Importer field to update")
	}
	// Swap the Importer for the identity-aware variant.
	e.replace(r.offset(imp.Value.Pos()), r.offset(imp.Value.End()),
		fmt.Sprintf("pluginsdk.ImporterValidatingIdentity(%s)", idRef))
	// Add the Identity schema field just before the Importer field.
	e.insert(r.lineStartOffset(imp.Pos()),
		fmt.Sprintf("Identity: &schema.ResourceIdentity{\nSchemaFunc: pluginsdk.GenerateIdentitySchema(%s),\n},\n\n", idRef))

	// Set the identity data in Create.
	if create := r.funcs[r.CreateFunc]; create != nil {
		if err := r.injectUntypedCreateIdentity(e, create); err != nil {
			return err
		}
	}

	// Add the identity-test go:generate directive.
	props := opts.IdentityProperties
	if props == "" {
		props = "name,resource_group_name"
	}
	if opts.ResourceName != "" {
		r.addGoGenerateEdit(e, opts.ResourceName, props)
	}
	return nil
}

// untypedImporterField returns the Importer key/value in the resource literal.
func (r *Resource) untypedImporterField() *ast.KeyValueExpr {
	for _, el := range r.resourceLit.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "Importer" {
				return kv
			}
		}
	}
	return nil
}

// injectUntypedCreateIdentity inserts a SetResourceIdentityData call after the
// d.SetId(...) call in Create.
func (r *Resource) injectUntypedCreateIdentity(e *editor, create *ast.FuncDecl) error {
	var (
		setIDStmt ast.Stmt
		idVar     string
	)
	ast.Inspect(create, func(n ast.Node) bool {
		es, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}
		call, ok := es.X.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "SetId" {
			return true
		}
		if x, ok := sel.X.(*ast.Ident); !ok || x.Name != "d" {
			return true
		}
		// d.SetId(id.ID()) -> the id variable is the selector base of the arg.
		if argCall, ok := call.Args[0].(*ast.CallExpr); ok {
			if argSel, ok := argCall.Fun.(*ast.SelectorExpr); ok {
				if base, ok := argSel.X.(*ast.Ident); ok {
					idVar = base.Name
				}
			}
		}
		setIDStmt = es
		return false
	})
	if setIDStmt == nil {
		return fmt.Errorf("could not find a d.SetId(...) call in Create")
	}
	if idVar == "" {
		idVar = "id"
	}
	ref := idVar
	if !r.idVarIsPointer(create, idVar) {
		ref = "&" + idVar
	}
	e.insert(r.offset(setIDStmt.End()),
		fmt.Sprintf("\nif err := pluginsdk.SetResourceIdentityData(d, %s); err != nil {\nreturn err\n}\n", ref))
	return nil
}

// idVarIsPointer reports whether idVar was bound from a Parse*ID call (a
// pointer) rather than a New*ID call (a value).
func (r *Resource) idVarIsPointer(fn *ast.FuncDecl, idVar string) bool {
	pointer := false
	ast.Inspect(fn, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || len(as.Lhs) == 0 {
			return true
		}
		lhs, ok := as.Lhs[0].(*ast.Ident)
		if !ok || lhs.Name != idVar || len(as.Rhs) != 1 {
			return true
		}
		if call, ok := as.Rhs[0].(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if strings.HasPrefix(sel.Sel.Name, "Parse") {
					pointer = true
					return false
				}
			}
		}
		return true
	})
	return pointer
}
