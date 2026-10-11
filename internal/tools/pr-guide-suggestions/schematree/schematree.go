// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package schematree walks Terraform Plugin SDK schema maps
// (map[string]*schema.Schema literals) purely at the syntax level and exposes a
// rooted tree of property Nodes for lint rules to inspect.
//
// Unlike the JSON-schema approach, nothing here compiles the provider or renders
// its schema: it only parses Go source, so it is fast and can run on a single
// file. It depends only on the standard go/ast; no type information is required.
package schematree

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/astutil"
)

// Result is the set of property nodes discovered in a package.
type Result struct {
	// All is every property node, in pre-order (parents before children).
	All []*Node
}

// Node is a single property: a key/value pair inside a map[string]*schema.Schema.
type Node struct {
	// Name is the leaf property name (the map key).
	Name string
	// Path is the dotted path from the resource root (e.g. "sku.name").
	Path string
	// Key is the map key expression; findings are anchored here so that
	// line-based diff filtering matches the property that changed.
	Key ast.Expr
	// Value is the map value expression.
	Value ast.Expr
	// Schema is the syntactic view of the schema literal, or nil when the value
	// is produced by a helper call (e.g. commonschema.Location()) and is opaque.
	Schema *SchemaLit
	// Parent is the enclosing block node, or nil at the resource root.
	Parent *Node
	// Siblings maps every property name at this level to its node (including
	// this one).
	Siblings map[string]*Node
	// Children maps nested block property names to their nodes (nil for
	// non-blocks).
	Children map[string]*Node
	// TopLevel reports whether this node is a direct resource/data source
	// argument or attribute.
	TopLevel bool
	// Resource is a best-effort resource/data source type name (e.g.
	// "azurerm_resource_group"); it may be empty when it cannot be determined.
	Resource string
	// Kind is "resource" or "data source", inferred from the file name.
	Kind string
}

// Build walks the parsed files and returns the property tree they define.
func Build(fset *token.FileSet, files []*ast.File) *Result {
	var schemaMaps []*ast.CompositeLit
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if cl, ok := n.(*ast.CompositeLit); ok && isSchemaMap(cl) {
				schemaMaps = append(schemaMaps, cl)
			}
			return true
		})
	}

	b := &builder{
		fset:      fset,
		files:     files,
		result:    &Result{},
		resources: resourceNames(files),
		funcs:     buildFuncIndex(files),
		walked:    map[*ast.CompositeLit]bool{},
		stack:     map[*ast.CompositeLit]bool{},
	}

	// A schema map is inline-"nested" when it is the direct Elem of a block in
	// another schema map; it is walked as a child of that block rather than as a
	// root.
	nested := make(map[*ast.CompositeLit]bool)
	for _, m := range schemaMaps {
		for _, kv := range mapEntries(m) {
			if cl, ok := unwrapSchemaLit(kv.Value); ok {
				if child, ok := newSchemaLit(cl).ElemResourceSchemaMap(); ok {
					nested[child] = true
				}
			}
		}
	}

	// Primary roots are maps that are neither inline-nested nor owned by a schema
	// helper function (those are reached by resolving the helper call).
	for _, m := range schemaMaps {
		if nested[m] || b.funcs.owned[m] {
			continue
		}
		b.walkRoot(m)
	}
	// Fallback: lint any remaining map (for example a helper whose caller is not
	// in the analysed files) so no schema is skipped.
	for _, m := range schemaMaps {
		if nested[m] || b.walked[m] {
			continue
		}
		b.walkRoot(m)
	}

	return b.result
}

type builder struct {
	fset      *token.FileSet
	files     []*ast.File
	result    *Result
	resources fileResourcesMap
	funcs     *funcIndex
	// walked records every map already visited, so a map consumed as a resolved
	// child is not also linted as a fallback root.
	walked map[*ast.CompositeLit]bool
	// stack holds the maps currently on the recursion path, to break cycles.
	stack map[*ast.CompositeLit]bool
}

func (b *builder) walkRoot(m *ast.CompositeLit) {
	if b.walked[m] {
		return
	}
	resource, kind := b.rootInfo(m)
	b.walk(m, nil, "", true, resource, kind)
}

// walk builds the property nodes for a schema map, resolving local schema-helper
// calls, recursing into nested blocks and recording the dotted path.
func (b *builder) walk(mapLit *ast.CompositeLit, parent *Node, prefix string, topLevel bool, resource, kind string) map[string]*Node {
	b.walked[mapLit] = true
	b.stack[mapLit] = true
	defer delete(b.stack, mapLit)

	level := make(map[string]*Node)
	var order []*Node

	for _, kv := range mapEntries(mapLit) {
		name, ok := stringLitValue(kv.Key)
		if !ok {
			continue
		}
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		node := &Node{
			Name:     name,
			Path:     path,
			Key:      kv.Key,
			Value:    kv.Value,
			Parent:   parent,
			TopLevel: topLevel,
			Resource: resource,
			Kind:     kind,
		}
		if cl, ok := b.resolveSchemaLit(kv.Value); ok {
			node.Schema = newSchemaLit(cl)
		}

		level[name] = node
		order = append(order, node)
		b.result.All = append(b.result.All, node)
	}

	for _, n := range order {
		n.Siblings = level
	}

	for _, n := range order {
		if n.Schema == nil {
			continue
		}
		childMap, ok := b.blockChildMap(n.Schema)
		if !ok || b.stack[childMap] {
			continue
		}
		n.Children = b.walk(childMap, n, n.Path, false, resource, kind)
	}

	return level
}

// rootInfo returns a best-effort resource type name and kind for a root schema
// map, using any ResourceType() method declared on the enclosing receiver.
func (b *builder) rootInfo(mapLit *ast.CompositeLit) (resource, kind string) {
	pos := mapLit.Pos()
	file := fileContaining(b.files, pos)
	if file == nil {
		return "", "resource"
	}

	kind = "resource"
	if name := b.fset.File(pos).Name(); strings.Contains(name, "data_source") || strings.Contains(name, "datasource") {
		kind = "data source"
	}

	fr := b.resources[file]
	if recv := enclosingReceiver(file, pos); recv != "" {
		if name, ok := fr.byReceiver[recv]; ok {
			return name, kind
		}
	}
	if len(fr.byReceiver) == 1 {
		for _, name := range fr.byReceiver {
			return name, kind
		}
	}

	return "", kind
}

// --- local schema-helper resolution ------------------------------------------

// funcIndex maps package-level schema-returning helper functions to the
// composite literal they return, so calls such as SchemaDefaultNodePool() and
// resourceKubernetesClusterNodePoolSchema() can be followed to their schema.
type funcIndex struct {
	schemaLit map[string]*ast.CompositeLit // func name -> returned *Schema literal
	mapLit    map[string]*ast.CompositeLit // func name -> returned schema map literal
	owned     map[*ast.CompositeLit]bool   // maps a helper exposes (its map, or its *Schema Elem map)
}

func buildFuncIndex(files []*ast.File) *funcIndex {
	idx := &funcIndex{
		schemaLit: map[string]*ast.CompositeLit{},
		mapLit:    map[string]*ast.CompositeLit{},
		owned:     map[*ast.CompositeLit]bool{},
	}

	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			lit := funcReturnCompositeLit(fn)
			if lit == nil {
				continue
			}
			switch {
			case isSchemaMap(lit):
				idx.mapLit[fn.Name.Name] = lit
			case isSelectorNamed(lit.Type, "Schema"):
				idx.schemaLit[fn.Name.Name] = lit
			}
		}
	}

	for _, m := range idx.mapLit {
		idx.owned[m] = true
	}
	for _, s := range idx.schemaLit {
		if m, ok := idx.schemaElemMap(s); ok {
			idx.owned[m] = true
		}
	}
	return idx
}

// schemaElemMap returns the nested schema map exposed by a *Schema literal whose
// Elem is a *Resource, resolving a helper call or IIFE for the Schema field.
func (idx *funcIndex) schemaElemMap(schemaLit *ast.CompositeLit) (*ast.CompositeLit, bool) {
	elem, ok := unwrapCompositeLit(compositeField(schemaLit, "Elem"))
	if !ok || !isSelectorNamed(elem.Type, "Resource") {
		return nil, false
	}
	sv := compositeField(elem, "Schema")
	if cl, ok := sv.(*ast.CompositeLit); ok && isSchemaMap(cl) {
		return cl, true
	}
	if name, ok := localCallName(sv); ok {
		if cl, ok := idx.mapLit[name]; ok {
			return cl, true
		}
	}
	if cl, ok := iifeCompositeLit(sv); ok && isSchemaMap(cl) {
		return cl, true
	}
	return nil, false
}

// resolveSchemaLit resolves a property value to its schema composite literal: a
// direct literal, a local schema-returning helper call, or an IIFE.
func (b *builder) resolveSchemaLit(e ast.Expr) (*ast.CompositeLit, bool) {
	if cl, ok := unwrapCompositeLit(e); ok {
		return cl, true
	}
	if name, ok := localCallName(e); ok {
		if cl, ok := b.funcs.schemaLit[name]; ok {
			return cl, true
		}
	}
	if cl, ok := iifeCompositeLit(e); ok && isSelectorNamed(cl.Type, "Schema") {
		return cl, true
	}
	return nil, false
}

// blockChildMap returns the nested schema map for a block schema, resolving the
// Elem's Schema field when it is a helper call or IIFE.
func (b *builder) blockChildMap(s *SchemaLit) (*ast.CompositeLit, bool) {
	return b.funcs.schemaElemMap(s.Lit)
}

// funcReturnCompositeLit returns the composite literal a function returns.
func funcReturnCompositeLit(fn *ast.FuncDecl) *ast.CompositeLit {
	if fn.Body == nil {
		return nil
	}
	return bodyReturnCompositeLit(fn.Body)
}

// bodyReturnCompositeLit returns the composite literal a body's last top-level
// return yields, following a single `return <ident>` back to that ident's
// assignment (the common `s := <lit>; ...; return s` helper pattern).
func bodyReturnCompositeLit(body *ast.BlockStmt) *ast.CompositeLit {
	ret := lastReturn(body)
	if ret == nil || len(ret.Results) != 1 {
		return nil
	}
	if cl, ok := unwrapCompositeLit(ret.Results[0]); ok {
		return cl
	}
	if id, ok := ret.Results[0].(*ast.Ident); ok {
		if rhs := lastAssignment(body, id.Name); rhs != nil {
			if cl, ok := unwrapCompositeLit(rhs); ok {
				return cl
			}
		}
	}
	return nil
}

// iifeCompositeLit returns the composite literal returned by an immediately
// invoked function literal, e.g. func() map[string]*schema.Schema { ... }().
func iifeCompositeLit(e ast.Expr) (*ast.CompositeLit, bool) {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	fl, ok := call.Fun.(*ast.FuncLit)
	if !ok || fl.Body == nil {
		return nil, false
	}
	cl := bodyReturnCompositeLit(fl.Body)
	return cl, cl != nil
}

// localCallName returns the identifier name of a call to a package-level
// function (e.g. SchemaDefaultNodePool()), and false for method or package
// qualified calls.
func localCallName(e ast.Expr) (string, bool) {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	if id, ok := call.Fun.(*ast.Ident); ok {
		return id.Name, true
	}
	return "", false
}

func lastReturn(body *ast.BlockStmt) *ast.ReturnStmt {
	var ret *ast.ReturnStmt
	for _, stmt := range body.List {
		if r, ok := stmt.(*ast.ReturnStmt); ok {
			ret = r
		}
	}
	return ret
}

func lastAssignment(body *ast.BlockStmt, name string) ast.Expr {
	var rhs ast.Expr
	for _, stmt := range body.List {
		as, ok := stmt.(*ast.AssignStmt)
		if !ok || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			continue
		}
		if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name == name {
			rhs = as.Rhs[0]
		}
	}
	return rhs
}

func compositeField(cl *ast.CompositeLit, name string) ast.Expr {
	if cl == nil {
		return nil
	}
	for _, kv := range fieldEntries(cl) {
		if identName(kv.Key) == name {
			return kv.Value
		}
	}
	return nil
}

// SchemaLit is a syntactic view over a schema.Schema composite literal.
type SchemaLit struct {
	Lit    *ast.CompositeLit
	fields map[string]*ast.KeyValueExpr
}

func newSchemaLit(cl *ast.CompositeLit) *SchemaLit {
	return &SchemaLit{Lit: cl, fields: astutil.CompositeLitFields(cl)}
}

// Declares reports whether the schema sets the named field.
func (s *SchemaLit) Declares(field string) bool { return s.fields[field] != nil }

// FieldValue returns the expression assigned to the named field, or nil.
func (s *SchemaLit) FieldValue(field string) ast.Expr {
	if kv := s.fields[field]; kv != nil {
		return kv.Value
	}
	return nil
}

// Bool returns the boolean value of a field, or false when unset/unreadable.
func (s *SchemaLit) Bool(field string) bool {
	if kv := s.fields[field]; kv != nil {
		if v := astutil.ExprBoolValue(kv.Value); v != nil {
			return *v
		}
	}
	return false
}

// Int returns the integer value of a field, or 0 when unset/unreadable.
func (s *SchemaLit) Int(field string) int {
	if kv := s.fields[field]; kv != nil {
		if v := astutil.ExprIntValue(kv.Value); v != nil {
			return *v
		}
	}
	return 0
}

// String returns the string value of a field, or "" when unset/unreadable.
func (s *SchemaLit) String(field string) string {
	if kv := s.fields[field]; kv != nil {
		if v := astutil.ExprStringValue(kv.Value); v != nil {
			return *v
		}
	}
	return ""
}

// ValueType returns the schema Type as its identifier name (e.g. "TypeString"),
// or "" when the Type field is absent or not a simple selector/identifier.
func (s *SchemaLit) ValueType() string {
	kv := s.fields["Type"]
	if kv == nil {
		return ""
	}
	switch v := kv.Value.(type) {
	case *ast.SelectorExpr:
		return v.Sel.Name
	case *ast.Ident:
		return v.Name
	}
	return ""
}

// IsCollection reports whether the schema Type is TypeList or TypeSet.
func (s *SchemaLit) IsCollection() bool {
	switch s.ValueType() {
	case "TypeList", "TypeSet":
		return true
	}
	return false
}

// ElemResourceSchemaMap returns the nested schema map when Elem is a
// *schema.Resource, and true when the property is a nested block.
func (s *SchemaLit) ElemResourceSchemaMap() (*ast.CompositeLit, bool) {
	cl, ok := unwrapCompositeLit(s.FieldValue("Elem"))
	if !ok || !isSelectorNamed(cl.Type, "Resource") {
		return nil, false
	}
	for _, kv := range fieldEntries(cl) {
		if identName(kv.Key) == "Schema" {
			if inner, ok := kv.Value.(*ast.CompositeLit); ok && isSchemaMap(inner) {
				return inner, true
			}
		}
	}
	return nil, false
}

// ElemIsScalarSchema reports whether Elem is a *schema.Schema, i.e. the property
// is a scalar array rather than a nested block.
func (s *SchemaLit) ElemIsScalarSchema() bool {
	cl, ok := unwrapCompositeLit(s.FieldValue("Elem"))
	return ok && isSelectorNamed(cl.Type, "Schema")
}

// --- syntactic helpers -------------------------------------------------------

// isSchemaMap reports whether a composite literal is a map[string]*schema.Schema.
func isSchemaMap(cl *ast.CompositeLit) bool {
	mt, ok := cl.Type.(*ast.MapType)
	if !ok || !astutil.IsStringType(mt.Key) {
		return false
	}
	v := mt.Value
	if star, ok := v.(*ast.StarExpr); ok {
		v = star.X
	}
	return isSelectorNamed(v, "Schema")
}

// unwrapSchemaLit returns the composite literal for a schema property value,
// unwrapping a leading &. It returns false for non-literal values such as helper
// calls.
func unwrapSchemaLit(e ast.Expr) (*ast.CompositeLit, bool) {
	return unwrapCompositeLit(e)
}

func unwrapCompositeLit(e ast.Expr) (*ast.CompositeLit, bool) {
	switch v := e.(type) {
	case *ast.CompositeLit:
		return v, true
	case *ast.UnaryExpr:
		if v.Op == token.AND {
			if cl, ok := v.X.(*ast.CompositeLit); ok {
				return cl, true
			}
		}
	}
	return nil, false
}

// mapEntries returns the key/value entries of a map composite literal.
func mapEntries(cl *ast.CompositeLit) []*ast.KeyValueExpr {
	out := make([]*ast.KeyValueExpr, 0, len(cl.Elts))
	for _, elt := range cl.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			out = append(out, kv)
		}
	}
	return out
}

// fieldEntries returns the field entries of a struct composite literal.
func fieldEntries(cl *ast.CompositeLit) []*ast.KeyValueExpr {
	return mapEntries(cl)
}

func isSelectorNamed(e ast.Expr, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == name
}

func identName(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

func stringLitValue(e ast.Expr) (string, bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	v, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return v, true
}

// --- resource-name discovery -------------------------------------------------

type fileResources struct {
	byReceiver map[string]string
}

type fileResourcesMap map[*ast.File]fileResources

// resourceNames maps each file to the resource type names declared by its
// ResourceType() string methods, keyed by receiver type name.
func resourceNames(files []*ast.File) fileResourcesMap {
	out := make(fileResourcesMap)
	for _, file := range files {
		fr := fileResources{byReceiver: map[string]string{}}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "ResourceType" || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			name, ok := singleReturnString(fn)
			if !ok {
				continue
			}
			fr.byReceiver[receiverTypeName(fn.Recv.List[0].Type)] = name
		}
		out[file] = fr
	}
	return out
}

func singleReturnString(fn *ast.FuncDecl) (string, bool) {
	if fn.Body == nil {
		return "", false
	}
	for _, stmt := range fn.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			continue
		}
		return stringLitValue(ret.Results[0])
	}
	return "", false
}

func receiverTypeName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.StarExpr:
		return receiverTypeName(v.X)
	case *ast.Ident:
		return v.Name
	case *ast.IndexExpr: // generic receiver
		return receiverTypeName(v.X)
	}
	return ""
}

func fileContaining(files []*ast.File, pos token.Pos) *ast.File {
	for _, file := range files {
		if file.Pos() <= pos && pos <= file.End() {
			return file
		}
	}
	return nil
}

func enclosingReceiver(file *ast.File, pos token.Pos) string {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}
		if fn.Body.Pos() <= pos && pos <= fn.Body.End() {
			return receiverTypeName(fn.Recv.List[0].Type)
		}
	}
	return ""
}
