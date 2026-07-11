package list_upgrade

import (
	"go/ast"
	"go/token"
	"strings"
)

// detectUntypedParent identifies whether the resource is a child listed under a
// parent (e.g. a virtual hub IP configuration under a virtual hub). It first
// tries the robust signal — the ID variable whose fields construct the child's
// New<Child>ID(...) — then falls back to matching Parse<X>ID(d.Get("<attr>"))
// candidates by prefix or uniqueness.
func (r *Resource) detectUntypedParent() {
	if r.detectParentFromConstruction() {
		return
	}

	type candidate struct{ pkg, base, attr string }
	var cands []candidate
	seen := map[string]bool{}
	for _, name := range []string{r.CreateFunc, r.ReadFunc, r.UpdateFunc} {
		fn := r.funcs[name]
		if fn == nil {
			continue
		}
		for _, c := range parentCandidatesInFunc(fn) {
			if c.base == r.IDBase || seen[c.attr] {
				continue
			}
			seen[c.attr] = true
			cands = append(cands, candidate{c.pkg, c.base, c.attr})
		}
	}
	if len(cands) == 0 {
		return
	}

	var chosen *candidate
	for i := range cands {
		if strings.HasPrefix(r.IDBase, cands[i].base) {
			chosen = &cands[i]
			break
		}
	}
	if chosen == nil && len(cands) == 1 {
		chosen = &cands[0]
	}
	if chosen == nil {
		return
	}
	r.setParent(chosen.pkg, chosen.base, chosen.attr)
}

// detectParentFromConstruction identifies the parent as the ID variable whose
// fields are passed to the child's New<Child>ID(...) constructor. This is robust
// even when the parent ID base is not a prefix of the child ID base and other
// unrelated IDs (e.g. remote_virtual_network_id) are parsed in the same func.
func (r *Resource) detectParentFromConstruction() bool {
	if r.IDBase == "" {
		return false
	}
	newFn := "New" + r.IDBase + "ID"
	var parentVar string
	for _, name := range []string{r.CreateFunc, r.UpdateFunc, r.ReadFunc} {
		if fn := r.funcs[name]; fn != nil {
			if v := newIDParentVar(fn, newFn); v != "" {
				parentVar = v
				break
			}
		}
	}
	if parentVar == "" {
		return false
	}
	for _, name := range []string{r.CreateFunc, r.UpdateFunc, r.ReadFunc} {
		if fn := r.funcs[name]; fn != nil {
			if pkg, base, attr := r.parentFromVar(fn, parentVar); base != "" && base != r.IDBase {
				r.setParent(pkg, base, attr)
				return true
			}
		}
	}
	return false
}

// detectTypedParent detects the parent scope for a typed resource from its
// Create/Update/Read closures. Typed resources parse the parent from a decoded
// model field (e.g. model.VirtualHubId) rather than d.Get, and their CRUD
// bodies live inside the inner sdk.ResourceFunc closure.
func (r *Resource) detectTypedParent() {
	if r.IDBase == "" {
		return
	}
	newFn := "New" + r.IDBase + "ID"
	var parentVar string
	for _, name := range []string{"Create", "Update", "Read"} {
		fd := r.methods[name]
		if fd == nil {
			continue
		}
		body := resourceFuncBody(fd)
		if body == nil {
			continue
		}
		if v := newIDParentVar(body, newFn); v != "" {
			parentVar = v
			break
		}
	}
	if parentVar == "" {
		return
	}
	for _, name := range []string{"Create", "Update", "Read"} {
		fd := r.methods[name]
		if fd == nil {
			continue
		}
		body := resourceFuncBody(fd)
		if body == nil {
			continue
		}
		if pkg, base, attr := r.parentFromVar(body, parentVar); base != "" && base != r.IDBase {
			r.setParent(pkg, base, attr)
			return
		}
	}
}

// setParent records the detected parent scope and the default parent-scoped list
// method (derived from the Get method; may be overridden by the caller for
// resources whose List method pluralises irregularly).
func (r *Resource) setParent(pkg, base, attr string) {
	r.ParentPackage = pkg
	r.ParentIDBase = base
	r.ParentIDType = base + "Id"
	r.ParentParseFunc = "Parse" + base + "ID"
	r.ParentValidateFunc = "Validate" + base + "ID"
	r.ParentAttr = attr
	if imp, ok := r.imports[pkg]; ok {
		r.ParentImportPath = strings.Trim(imp.Path.Value, `"`)
	}
	r.ListMethod = strings.TrimSuffix(r.GetMethod, "Get") + "List"
}

// newIDParentVar finds a `<pkg>.<newFn>(...)` call and returns the variable whose
// fields are passed as arguments most often (the parent scope ID variable).
func newIDParentVar(scope ast.Node, newFn string) string {
	var out string
	ast.Inspect(scope, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != newFn {
			return true
		}
		counts := map[string]int{}
		var order []string
		for _, arg := range call.Args {
			if s, ok := arg.(*ast.SelectorExpr); ok {
				if id, ok := s.X.(*ast.Ident); ok {
					if counts[id.Name] == 0 {
						order = append(order, id.Name)
					}
					counts[id.Name]++
				}
			}
		}
		best, bestN := "", 0
		for _, name := range order {
			if counts[name] > bestN {
				best, bestN = name, counts[name]
			}
		}
		if best != "" {
			out = best
			return false
		}
		return true
	})
	return out
}

// parentFromVar finds `<varName>, err := <pkg>.Parse<Base>ID(<arg>)` within the
// given scope and returns the package, ID base and config attribute. The arg is
// either `d.Get("<attr>").(string)` (untyped) or `<model>.<Field>` (typed), in
// which case the attribute is resolved from the model field's tfschema tag.
func (r *Resource) parentFromVar(scope ast.Node, varName string) (pkg, base, attr string) {
	ast.Inspect(scope, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || len(as.Lhs) == 0 || len(as.Rhs) != 1 {
			return true
		}
		lhs, ok := as.Lhs[0].(*ast.Ident)
		if !ok || lhs.Name != varName {
			return true
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		fnName := sel.Sel.Name
		if !strings.HasPrefix(fnName, "Parse") || !strings.HasSuffix(fnName, "ID") {
			return true
		}
		a := r.resolveParentAttr(call.Args[0])
		if a == "" {
			return true
		}
		pkg = pkgIdent.Name
		base = strings.TrimSuffix(strings.TrimPrefix(fnName, "Parse"), "ID")
		attr = a
		return false
	})
	return pkg, base, attr
}

// resolveParentAttr returns the terraform config attribute for a parent parse
// argument: the literal key of `d.Get("<key>")` or the tfschema tag of the model
// field in `<model>.<Field>` (falling back to a snake_case of the field name).
func (r *Resource) resolveParentAttr(arg ast.Expr) string {
	if k := configGetKey(arg); k != "" {
		return k
	}
	if sel, ok := arg.(*ast.SelectorExpr); ok {
		if tag := r.modelFieldTFTag(sel.Sel.Name); tag != "" {
			return tag
		}
		return toSnakeCase(sel.Sel.Name)
	}
	return ""
}

// modelFieldTFTag returns the tfschema tag of the named field on the resource's
// model struct, or "" when it cannot be found.
func (r *Resource) modelFieldTFTag(fieldName string) string {
	if r.ModelStruct == "" {
		return ""
	}
	for _, decl := range r.file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != r.ModelStruct {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok || st.Fields == nil {
				continue
			}
			for _, f := range st.Fields.List {
				if len(f.Names) != 1 || f.Names[0].Name != fieldName || f.Tag == nil {
					continue
				}
				return structTagValue(strings.Trim(f.Tag.Value, "`"), "tfschema")
			}
		}
	}
	return ""
}

// structTagValue extracts the value of a struct tag key from a raw tag string.
func structTagValue(tag, key string) string {
	needle := key + `:"`
	idx := strings.Index(tag, needle)
	if idx == -1 {
		return ""
	}
	rest := tag[idx+len(needle):]
	if end := strings.Index(rest, `"`); end != -1 {
		return rest[:end]
	}
	return ""
}

// toSnakeCase converts a CamelCase identifier to snake_case (e.g. VirtualHubId ->
// virtual_hub_id). Used only as a fallback when a tfschema tag is unavailable.
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(ch + ('a' - 'A'))
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// parentCandidatesInFunc returns every (pkg, base, attr) from a
// `pkg.Parse<base>ID(d.Get("<attr>").(string))` call in the function.
func parentCandidatesInFunc(fn *ast.FuncDecl) []struct{ pkg, base, attr string } {
	var out []struct{ pkg, base, attr string }
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		fnName := sel.Sel.Name
		if !strings.HasPrefix(fnName, "Parse") || !strings.HasSuffix(fnName, "ID") {
			return true
		}
		attr := configGetKey(call.Args[0])
		if attr == "" {
			return true
		}
		out = append(out, struct{ pkg, base, attr string }{
			pkg:  pkgIdent.Name,
			base: strings.TrimSuffix(strings.TrimPrefix(fnName, "Parse"), "ID"),
			attr: attr,
		})
		return true
	})
	return out
}

// configGetKey returns the literal key from a `d.Get("<key>").(string)` expression.
func configGetKey(expr ast.Expr) string {
	ta, ok := expr.(*ast.TypeAssertExpr)
	if !ok {
		return ""
	}
	getCall, ok := ta.X.(*ast.CallExpr)
	if !ok || len(getCall.Args) != 1 {
		return ""
	}
	getSel, ok := getCall.Fun.(*ast.SelectorExpr)
	if !ok || getSel.Sel.Name != "Get" {
		return ""
	}
	if x, ok := getSel.X.(*ast.Ident); !ok || x.Name != "d" {
		return ""
	}
	lit, ok := getCall.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, `"`)
}
