// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes

import (
	"go/ast"
	"strings"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const AZNR003Doc = `check that expand/flatten functions are defined as receiver methods

The AZNR003 analyzer reports when expand* or flatten* functions are defined as 
global/package-level functions instead of receiver methods on a resource type.

This check only applies to typed resources.

Example violation:

	// Global function - should be a receiver method
	func expandCustomerManagedKey(input []CustomerManagedKey) (*Encryption, error) {
	    // ...
	}

	func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	    // ...
	}

Correct usage:

	func (r AIServices) expandCustomerManagedKey(input []CustomerManagedKey) (*Encryption, error) {
	    // ...
	}

	func (r AIServices) flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs {
	    // ...
	}
`

const aznr003Name = "AZNR003"

var AZNR003Analyzer = &analysis.Analyzer{
	Name:     aznr003Name,
	Doc:      AZNR003Doc,
	Run:      runAZNR003,
	Requires: []*analysis.Analyzer{inspect.Analyzer, commentignore.Analyzer, schema.TypedResourceInfoAnalyzer},
}

func runAZNR003(pass *analysis.Pass) (interface{}, error) {
	if helper.ShouldSkipPackageForResourceAnalysis(pass.Pkg.Path()) {
		return nil, nil
	}

	ignorer, ok := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	if !ok {
		return nil, nil
	}

	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Get typed resources from analyzer
	typedResources, ok := pass.ResultOf[schema.TypedResourceInfoAnalyzer].([]*helper.TypedResourceInfo)
	if !ok || len(typedResources) == 0 {
		return nil, nil
	}

	// Build a map of filenames that contain typed resources
	typedResourceFiles := make(map[string]bool)
	for _, resource := range typedResources {
		if resource.ArgumentsFunc != nil {
			filename := pass.Fset.Position(resource.ArgumentsFunc.Pos()).Filename
			typedResourceFiles[filename] = true
		}
	}

	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil {
			return
		}

		// Only check files that contain typed resources
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		if !typedResourceFiles[filename] {
			return
		}

		// Check git filter
		pos := pass.Fset.Position(funcDecl.Pos())
		if !loader.ShouldReport(pos.Filename, pos.Line) {
			return
		}

		// Check if function name starts with "expand" or "flatten" (case-insensitive for first char after prefix)
		funcName := funcDecl.Name.Name
		isExpandFunc := strings.HasPrefix(strings.ToLower(funcName), "expand")
		isFlattenFunc := strings.HasPrefix(strings.ToLower(funcName), "flatten")
		if !isExpandFunc && !isFlattenFunc {
			return
		}

		// Check if it's already a receiver method
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			// Already a receiver method, no issue
			return
		}

		// Check comment ignore
		if ignorer.ShouldIgnore(aznr003Name, funcDecl) {
			return
		}

		// It's a global expand/flatten function - report issue
		funcType := "expand"
		if isFlattenFunc {
			funcType = "flatten"
		}

		pass.Reportf(funcDecl.Name.Pos(), "%s: %s function '%s' should be defined as a receiver method on the resource type, not as a global function\n",
			aznr003Name,
			funcType,
			helper.IssueLine(funcName))
	})

	return nil, nil
}
