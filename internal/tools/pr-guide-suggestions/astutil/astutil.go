// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package astutil provides the small set of go/ast helpers the linter needs to
// read schema composite literals. They are intentionally minimal and depend only
// on the standard library.
package astutil

import (
	"go/ast"
	"go/token"
	"strconv"
)

// CompositeLitFields returns the keyed fields of a struct composite literal,
// keyed by field name. Non-keyed or non-identifier-keyed elements are skipped.
func CompositeLitFields(cl *ast.CompositeLit) map[string]*ast.KeyValueExpr {
	out := make(map[string]*ast.KeyValueExpr, len(cl.Elts))
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if id, ok := kv.Key.(*ast.Ident); ok {
			out[id.Name] = kv
		}
	}
	return out
}

// ExprBoolValue returns the value of a bool literal (true/false), or nil.
func ExprBoolValue(e ast.Expr) *bool {
	id, ok := e.(*ast.Ident)
	if !ok {
		return nil
	}
	switch id.Name {
	case "true":
		v := true
		return &v
	case "false":
		v := false
		return &v
	}
	return nil
}

// ExprIntValue returns the value of an integer literal (including a negated
// one), or nil.
func ExprIntValue(e ast.Expr) *int {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind == token.INT {
			if n, err := strconv.Atoi(v.Value); err == nil {
				return &n
			}
		}
	case *ast.UnaryExpr:
		if v.Op == token.SUB {
			if inner := ExprIntValue(v.X); inner != nil {
				n := -*inner
				return &n
			}
		}
	}
	return nil
}

// ExprStringValue returns the unquoted value of a string literal, or nil.
func ExprStringValue(e ast.Expr) *string {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}
	if s, err := strconv.Unquote(lit.Value); err == nil {
		return &s
	}
	return nil
}

// IsStringType reports whether an expression is the identifier `string`.
func IsStringType(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "string"
}
