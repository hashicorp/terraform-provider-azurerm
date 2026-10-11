// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/astutil"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl004 = &Rule{
	ID:       "SL004",
	Name:     "avoid-none-value",
	Severity: Warning,
	Check:    checkSL004,
}

// noneValues are the sentinel enum members that should be omitted from a
// user-settable enum (Terraform null models "unset") and normalised in
// Create/Read instead. The order is used for deterministic output.
var noneValues = []string{"None", "Off", "Default", "Disabled"}

func checkSL004(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || !userSettable(n.Schema) {
			continue
		}
		vf := n.Schema.FieldValue(fieldValidateFunc)
		if vf == nil {
			continue
		}

		accepted := acceptedNoneValues(vf)
		if len(accepted) == 0 {
			continue
		}

		report(n, fmt.Sprintf("property %q accepts %s via validation; omit these values (use Terraform null) and normalise in Create/Read", n.Path, strings.Join(accepted, ", ")), nil)
	}
}

// acceptedNoneValues finds StringInSlice validators anywhere within the
// ValidateFunc expression (including those nested inside validation.Any/All
// wrappers) and returns the sentinel values they accept, in noneValues order. A
// sentinel is recognised either as a string literal ("None") or as an enum
// constant whose name ends in the sentinel (e.g. string(pkg.ChannelNone)).
func acceptedNoneValues(vf ast.Expr) []string {
	found := map[string]bool{}
	ast.Inspect(vf, func(nd ast.Node) bool {
		call, ok := nd.(*ast.CallExpr)
		if !ok || !isSelectorCall(call, "StringInSlice") || len(call.Args) == 0 {
			return true
		}
		collectEnumSentinels(call.Args[0], found)
		return true
	})

	var out []string
	for _, s := range noneValues {
		if found[s] {
			out = append(out, s)
		}
	}
	return out
}

// collectEnumSentinels records which sentinels appear in a []string{...}
// validator argument. A non-literal argument (e.g. PossibleValuesForX()) yields
// nothing.
func collectEnumSentinels(arg ast.Expr, found map[string]bool) {
	cl, ok := arg.(*ast.CompositeLit)
	if !ok {
		return
	}
	if at, ok := cl.Type.(*ast.ArrayType); !ok || !astutil.IsStringType(at.Elt) {
		return
	}
	for _, elt := range cl.Elts {
		tok, isLiteral := enumToken(elt)
		if tok == "" {
			continue
		}
		for _, s := range noneValues {
			if isLiteral {
				if tok == s {
					found[s] = true
				}
			} else if strings.HasSuffix(tok, s) {
				found[s] = true
			}
		}
	}
}

// enumToken extracts a comparable token from a StringInSlice element: the string
// value for a literal, or the constant identifier name for an enum reference
// such as pkg.ChannelNone or string(pkg.ChannelNone).
func enumToken(e ast.Expr) (tok string, isLiteral bool) {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			if s, err := strconv.Unquote(v.Value); err == nil {
				return s, true
			}
		}
	case *ast.CallExpr: // string(pkg.Const)
		if id, ok := v.Fun.(*ast.Ident); ok && id.Name == "string" && len(v.Args) == 1 {
			return selectorOrIdentName(v.Args[0]), false
		}
	case *ast.SelectorExpr:
		return v.Sel.Name, false
	case *ast.Ident:
		return v.Name, false
	}
	return "", false
}

func selectorOrIdentName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.SelectorExpr:
		return v.Sel.Name
	case *ast.Ident:
		return v.Name
	}
	return ""
}
