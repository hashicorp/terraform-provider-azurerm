package passes

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const AZRE001Doc = `check for fixed error strings using fmt.Errorf instead of errors.New

The AZRE001 analyzer reports cases where fixed error strings (without format placeholders)
use fmt.Errorf() instead of errors.New().

Example violations:
  fmt.Errorf("something went wrong")  // should use errors.New()
  
Valid usage:
  errors.New("something went wrong")
  fmt.Errorf("value %s is invalid", value)  // has placeholder, OK`

const azre001Name = "AZRE001"

var AZRE001Analyzer = &analysis.Analyzer{
	Name:     azre001Name,
	Doc:      AZRE001Doc,
	Run:      runAZRE001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runAZRE001(pass *analysis.Pass) (interface{}, error) {
	// Skip migration packages
	if strings.Contains(pass.Pkg.Path(), "/migration") {
		return nil, nil
	}

	// Pre-filter: Build set of changed files that import "fmt"
	relevantFiles := make(map[string]bool)
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		// Skip if not changed
		if !loader.IsFileChanged(filename) {
			continue
		}

		// Skip test files
		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		// Check if file imports "fmt"
		importsFmt := false
		for _, imp := range f.Imports {
			if imp.Path.Value == `"fmt"` {
				importsFmt = true
				break
			}
		}

		if importsFmt {
			relevantFiles[filename] = true
		}
	}

	// Early return if no relevant files
	if len(relevantFiles) == 0 {
		return nil, nil
	}

	// Get the shared inspector
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Pre-filter: only look at CallExpr nodes
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		// Use pre-computed Position for both filename and line number
		pos := pass.Fset.Position(call.Pos())
		filename := pos.Filename
		if !relevantFiles[filename] {
			return
		}

		// Check if it's a selector expression (pkg.Function)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// Check if it's calling Errorf
		if sel.Sel.Name != "Errorf" {
			return
		}

		// Check if the package is fmt
		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "fmt" {
			return
		}

		// Check if there are arguments
		if len(call.Args) != 1 {
			return
		}

		// Check if the first argument is a string literal
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		// Get the string value
		formatStr := lit.Value

		// Check if the string value contains any placeholders (%v, %s, %d, %+v, etc.)
		// If it doesn't contain %, it's a fixed string and should use errors.New()
		if !strings.Contains(formatStr, "%") {
			// Reuse pos from earlier to avoid duplicate Position lookup
			if loader.ShouldReport(filename, pos.Line) {
				pass.Reportf(call.Pos(), "%s: fixed error strings should use %s instead of %s\n",
					azre001Name,
					helper.FixedCode("errors.New()"),
					helper.IssueLine("fmt.Errorf()"))
			}
		}
	})

	return nil, nil
}
