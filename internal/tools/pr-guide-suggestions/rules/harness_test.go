// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

// finding is a rule diagnostic resolved to its property path for assertions.
type finding struct {
	rule   string
	path   string
	msg    string
	fix    string
	rename string
}

// wrap embeds schema-map entries in a minimal native resource so the source
// parses. Because the linter is purely syntactic, the source need not
// type-check or resolve imports.
func wrap(entries string) string {
	return `package x
func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
` + entries + `
		},
	}
}
`
}

// runRule parses src, builds the schema tree, and runs a single rule's check
// over it, returning each finding resolved to its property path. Only the one
// rule runs, so assertions see exactly its output.
func runRule(t *testing.T, rule *Rule, src string) []finding {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "x_resource.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	res := schematree.Build(fset, []*ast.File{file})

	var out []finding
	rule.Check(res, func(n *schematree.Node, message string, fix *Fix) {
		f := finding{rule: rule.ID, path: n.Path, msg: message}
		if fix != nil {
			f.fix = fix.Suggestion
			f.rename = fix.Rename
		}
		out = append(out, f)
	})
	return out
}

// flagged reports whether the given property path was flagged.
func flagged(fs []finding, path string) bool {
	for _, f := range fs {
		if f.path == path {
			return true
		}
	}
	return false
}

// at returns the finding for a property path, or nil.
func at(fs []finding, path string) *finding {
	for i := range fs {
		if fs[i].path == path {
			return &fs[i]
		}
	}
	return nil
}
