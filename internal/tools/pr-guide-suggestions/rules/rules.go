// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package rules implements the pr-guide-suggestions checks (SL001-SL013) as
// plain checks over the syntactic schema tree produced by the schematree
// package.
//
// Each rule is a *Rule bundling its stable ID, human name, severity and a Check
// function. A check reports findings against schema nodes; the driver anchors
// each finding at the property's map key so that line-based diff filtering
// matches the changed property. Fixable rules pass a *Fix; a Fix carrying a
// Rename is applied in place by -fix, others are only printed as suggestions.
package rules

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/astutil"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

// Severity classifies a finding for reporting and the process exit code.
type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

// Fix describes how to remediate a finding. Suggestion is a human-readable
// summary printed with each finding. When Rename is non-empty the fix is
// auto-applicable: -fix renames the property by replacing every quoted
// occurrence of its current name with Rename throughout the file (its schema
// key, d.Get/d.Set calls, tfschema tags and cross-field references).
type Fix struct {
	Suggestion string
	Rename     string
}

// ReportFunc records a rule finding against a schema node, with an optional fix.
type ReportFunc func(n *schematree.Node, message string, fix *Fix)

// Rule bundles a check's metadata with the function that runs it.
type Rule struct {
	ID       string
	Name     string
	Severity Severity
	Fixable  bool
	Check    func(*schematree.Result, ReportFunc)
}

// Registry lists every rule the linter ships, in ID order.
var Registry = []*Rule{
	sl001, sl002, sl003, sl004, sl005, sl006, sl007,
	sl008, sl009, sl010, sl011, sl012, sl013,
}

// ByID returns the rule with the given ID, or nil.
func ByID(id string) *Rule {
	for _, r := range Registry {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// schema.Schema field names.
const (
	fieldType         = "Type"
	fieldOptional     = "Optional"
	fieldRequired     = "Required"
	fieldDescription  = "Description"
	fieldMinItems     = "MinItems"
	fieldMaxItems     = "MaxItems"
	fieldValidateFunc = "ValidateFunc"
	fieldElem         = "Elem"
	fieldAtLeastOneOf = "AtLeastOneOf"
	fieldExactlyOneOf = "ExactlyOneOf"
	fieldValidateDiag = "ValidateDiagFunc"
)

// schema.ValueType names.
const (
	typeBool   = "TypeBool"
	typeInt    = "TypeInt"
	typeFloat  = "TypeFloat"
	typeString = "TypeString"
)

// userSettable reports whether a property is an argument (Optional or Required).
func userSettable(s *schematree.SchemaLit) bool {
	return s.Bool(fieldOptional) || s.Bool(fieldRequired)
}

// hasValidation reports whether a property sets ValidateFunc or ValidateDiagFunc.
func hasValidation(s *schematree.SchemaLit) bool {
	return s.Declares(fieldValidateFunc) || s.Declares(fieldValidateDiag)
}

// setNonZeroInt reports whether an int field is declared with a non-zero value
// (a non-literal value is assumed meaningful and therefore "set").
func setNonZeroInt(s *schematree.SchemaLit, field string) bool {
	if !s.Declares(field) {
		return false
	}
	if v := astutil.ExprIntValue(s.FieldValue(field)); v != nil {
		return *v != 0
	}
	return true
}

// emptyStringLiteral reports whether e is the empty string literal "".
func emptyStringLiteral(e ast.Expr) bool {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	v, err := strconv.Unquote(lit.Value)
	return err == nil && v == ""
}

// isSelectorCall reports whether a call is <anyPkg>.<funcName>(...).
func isSelectorCall(call *ast.CallExpr, funcName string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == funcName
}

// isSelectorFunc reports whether an expression is a bare <anyPkg>.<funcName>
// reference (not a call).
func isSelectorFunc(e ast.Expr, funcName string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == funcName
}
