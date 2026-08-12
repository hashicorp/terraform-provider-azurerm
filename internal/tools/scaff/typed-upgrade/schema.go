// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// FieldKind classifies the Go model type for a schema field.
type FieldKind int

const (
	FieldString  FieldKind = iota // string
	FieldInt                      // int64
	FieldFloat                    // float64
	FieldBool                     // bool
	FieldListStr                  // []string (list/set of TypeString primitives)
	FieldListObj                  // []BlockName (list/set of nested objects)
	FieldMap                      // map[string]string
	FieldUnknown                  // fall back to string
)

// GoType returns the Go model type string for this kind.
func (k FieldKind) GoType(blockName string) string {
	switch k {
	case FieldInt:
		return "int64"
	case FieldFloat:
		return "float64"
	case FieldBool:
		return "bool"
	case FieldListStr:
		return "[]string"
	case FieldListObj:
		if blockName != "" {
			return "[]" + blockName
		}
		return "[]interface{}"
	case FieldMap:
		return "map[string]string"
	default:
		return "string"
	}
}

// SchemaField captures the essential metadata for one entry in a
// pluginsdk.Resource.Schema map.
type SchemaField struct {
	TFName    string    // schema key, e.g. "maps_account_id"
	GoField   string    // PascalCase struct field, e.g. "MapsAccountId"
	Kind      FieldKind // Go model type classification
	ElemKind  FieldKind // for FieldListStr: element kind; unused for FieldListObj
	BlockName string    // for FieldListObj: PascalCase block model name
	MaxItems  int       // for FieldListObj: 1 = single nested block, 0 = list
	Required  bool
	Optional  bool
	Computed  bool
	ForceNew  bool
	// RawSchema is the raw Go source text of the schema value expression,
	// used verbatim in the Arguments/Attributes method bodies.
	RawSchema string
	// IsHelper is true when the value is a function call expression
	// (e.g. commonschema.Location()) rather than a schema literal.
	IsHelper bool
}

// IsArgument reports whether this field belongs in Arguments() (user-configurable).
func (f *SchemaField) IsArgument() bool {
	return f.Required || f.Optional || (!f.Computed && !f.Required && !f.Optional)
}

// IsAttribute reports whether this field belongs in Attributes() (computed-only).
func (f *SchemaField) IsAttribute() bool {
	return f.Computed && !f.Required && !f.Optional
}

// extractSchemaFields extracts the schema field definitions from the Schema
// entry of a pluginsdk.Resource{} composite literal.
func extractSchemaFields(src []byte, fset *token.FileSet, resourceLit *ast.CompositeLit) []*SchemaField {
	var schemaMapLit *ast.CompositeLit

	for _, el := range resourceLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Schema" {
			continue
		}
		// Schema value is map[string]*pluginsdk.Schema{...}
		if lit, ok := kv.Value.(*ast.CompositeLit); ok {
			schemaMapLit = lit
		}
		break
	}

	if schemaMapLit == nil {
		return nil
	}

	var fields []*SchemaField
	for _, el := range schemaMapLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		// Key is a string literal, e.g. "name"
		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}
		tfName := strings.Trim(keyLit.Value, `"`)
		rawSchema := nodeText(src, fset, kv.Value.Pos(), kv.Value.End())
		f := buildSchemaField(src, fset, tfName, rawSchema, kv.Value)
		fields = append(fields, f)
	}
	return fields
}

// buildSchemaField builds a SchemaField from a schema value expression.
func buildSchemaField(src []byte, fset *token.FileSet, tfName, rawSchema string, expr ast.Expr) *SchemaField {
	f := &SchemaField{
		TFName:    tfName,
		GoField:   snakeToCamel(tfName),
		RawSchema: rawSchema,
	}

	// Helper call: commonschema.Location(), commonschema.Tags(), etc.
	if call, ok := expr.(*ast.CallExpr); ok {
		f.IsHelper = true
		funcName := callFuncName(call)
		f.Kind, f.Required, f.Optional, f.Computed = classifyHelper(funcName)
		return f
	}

	// Composite literal: &pluginsdk.Schema{...} or pluginsdk.Schema{...}
	lit := unwrapSchemaLit(expr)
	if lit == nil {
		f.Kind = FieldUnknown
		f.Optional = true
		return f
	}

	var typeExpr, elemExpr ast.Expr
	for _, el := range lit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "Type":
			typeExpr = kv.Value
		case "Required":
			f.Required = isBoolTrue(kv.Value)
		case "Optional":
			f.Optional = isBoolTrue(kv.Value)
		case "Computed":
			f.Computed = isBoolTrue(kv.Value)
		case "ForceNew":
			f.ForceNew = isBoolTrue(kv.Value)
		case "Elem":
			elemExpr = kv.Value
		case "MaxItems":
			if basic, ok := kv.Value.(*ast.BasicLit); ok {
				fmt.Sscanf(basic.Value, "%d", &f.MaxItems)
			}
		}
	}

	f.Kind, f.ElemKind, f.BlockName = classifyType(src, fset, typeExpr, elemExpr, tfName)
	return f
}

// classifyType determines the FieldKind from the Type and Elem expressions.
func classifyType(src []byte, fset *token.FileSet, typeExpr, elemExpr ast.Expr, tfName string) (kind, elemKind FieldKind, blockName string) {
	typeName := selectorName(typeExpr)
	switch typeName {
	case "TypeString":
		return FieldString, FieldUnknown, ""
	case "TypeInt":
		return FieldInt, FieldUnknown, ""
	case "TypeFloat":
		return FieldFloat, FieldUnknown, ""
	case "TypeBool":
		return FieldBool, FieldUnknown, ""
	case "TypeMap":
		return FieldMap, FieldUnknown, ""
	case "TypeList", "TypeSet":
		if elemExpr == nil {
			return FieldListStr, FieldString, ""
		}
		return classifyListElem(src, fset, elemExpr, tfName)
	}
	return FieldUnknown, FieldUnknown, ""
}

// classifyListElem determines the element type for a list/set field.
func classifyListElem(src []byte, fset *token.FileSet, elemExpr ast.Expr, tfName string) (kind, elemKind FieldKind, blockName string) {
	// Unwrap & operator
	if unary, ok := elemExpr.(*ast.UnaryExpr); ok {
		elemExpr = unary.X
	}
	lit, ok := elemExpr.(*ast.CompositeLit)
	if !ok {
		return FieldListStr, FieldString, ""
	}
	typeName := compositeTypeName(lit)
	switch typeName {
	case "pluginsdk.Resource", "Resource":
		// Nested block: derive block name from tfName
		return FieldListObj, FieldUnknown, snakeToCamel(tfName) + "Model"
	case "pluginsdk.Schema", "Schema":
		// Inline schema element - check its Type
		for _, el := range lit.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Type" {
				continue
			}
			switch selectorName(kv.Value) {
			case "TypeString":
				return FieldListStr, FieldString, ""
			case "TypeInt":
				return FieldListStr, FieldInt, ""
			case "TypeBool":
				return FieldListStr, FieldBool, ""
			}
		}
		return FieldListStr, FieldString, ""
	}
	return FieldListStr, FieldString, ""
}

// classifyHelper maps a helper function name to field kind and cardinality.
func classifyHelper(funcName string) (kind FieldKind, required, optional, computed bool) {
	lower := strings.ToLower(funcName)
	switch {
	case strings.Contains(lower, "tags"):
		return FieldMap, false, true, true
	case strings.Contains(lower, "location"):
		return FieldString, false, true, false
	case strings.Contains(lower, "requiredsimplescope"):
		return FieldString, true, false, false
	default:
		return FieldString, false, true, false
	}
}

// callFuncName returns "Package.Func" or "Func" for a call expression.
func callFuncName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		if x, ok := fn.X.(*ast.Ident); ok {
			return x.Name + "." + fn.Sel.Name
		}
		return fn.Sel.Name
	case *ast.Ident:
		return fn.Name
	}
	return ""
}

// selectorName returns the selector field name (e.g. "TypeString" from "pluginsdk.TypeString").
func selectorName(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.Ident:
		return e.Name
	}
	return ""
}

// compositeTypeName returns the type name of a composite literal
// (e.g. "pluginsdk.Resource" from &pluginsdk.Resource{}).
func compositeTypeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		return typeExprName(e.Type)
	case *ast.UnaryExpr:
		if cl, ok := e.X.(*ast.CompositeLit); ok {
			return typeExprName(cl.Type)
		}
	}
	return ""
}

func typeExprName(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		if x, ok := e.X.(*ast.Ident); ok {
			return x.Name + "." + e.Sel.Name
		}
		return e.Sel.Name
	case *ast.Ident:
		return e.Name
	}
	return ""
}

// unwrapSchemaLit unwraps &pluginsdk.Schema{...} or pluginsdk.Schema{...}
// to the underlying composite literal, or returns nil.
func unwrapSchemaLit(expr ast.Expr) *ast.CompositeLit {
	if unary, ok := expr.(*ast.UnaryExpr); ok {
		if lit, ok := unary.X.(*ast.CompositeLit); ok {
			return lit
		}
	}
	if lit, ok := expr.(*ast.CompositeLit); ok {
		return lit
	}
	return nil
}

// isBoolTrue reports whether expr is the identifier "true".
func isBoolTrue(expr ast.Expr) bool {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name == "true"
	}
	return false
}

// snakeToCamel converts a snake_case string to PascalCase.
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// nodeText returns the raw source text between two token positions.
func nodeText(src []byte, fset *token.FileSet, start, end token.Pos) string {
	s := fset.Position(start).Offset
	e := fset.Position(end).Offset
	if s < 0 || e > len(src) || s > e {
		return ""
	}
	return string(src[s:e])
}
