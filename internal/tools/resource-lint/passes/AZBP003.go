// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes

import (
	"go/ast"
	"go/types"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/reporting"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	PointerPkgPath = "github.com/hashicorp/go-azure-helpers/lang/pointer"
	azbp003Name    = "AZBP003"
)

const AZBP003Doc = `check enum conversion for go-azure-sdk

The AZBP003 analyzer checks that when enum types are converted using pointer.To() with explicit type conversion instead of pointer.ToEnum[]

Example violation:
  return &managedclusters.ManagedClusterBootstrapProfile{
    ArtifactSource: pointer.To(managedclusters.ArtifactSource(config["artifact_source"].(string))),
  }

Valid usage:
  return &managedclusters.ManagedClusterBootstrapProfile{
    ArtifactSource: pointer.ToEnum[managedclusters.ArtifactSource](config["artifact_source"].(string)),
  }`

var AZBP003Analyzer = &analysis.Analyzer{
	Name:     azbp003Name,
	Doc:      AZBP003Doc,
	Run:      runAZBP003,
	Requires: []*analysis.Analyzer{inspect.Analyzer, commentignore.Analyzer},
}

func runAZBP003(pass *analysis.Pass) (interface{}, error) {
	ignorer, ok := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	if !ok {
		return nil, nil
	}

	if helper.ShouldSkipPackageForResourceAnalysis(pass.Pkg.Path()) {
		return nil, nil
	}

	relevantFiles := make(map[string]bool)
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		importsPointer := false
		for _, imp := range f.Imports {
			if imp.Path.Value == `"`+PointerPkgPath+`"` {
				importsPointer = true
				break
			}
		}

		if importsPointer {
			relevantFiles[filename] = true
		}
	}

	if len(relevantFiles) == 0 {
		return nil, nil
	}

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
		if !ok || len(call.Args) != 1 {
			return
		}

		// pattern: pointer.To(sdk.EnumConstant(xxx(string)))
		argCall, ok := call.Args[0].(*ast.CallExpr)
		if !ok {
			return
		}

		pos := pass.Fset.Position(call.Pos())
		filename := pos.Filename
		if !relevantFiles[filename] {
			return
		}

		// selector expressions: e.g. pointer.To
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "To" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "pointer" {
			return
		}

		// Validate pkg path
		obj := pass.TypesInfo.Uses[ident]
		if obj == nil {
			return
		}

		pkg, ok := obj.(*types.PkgName)
		if !ok || pkg.Imported().Path() != PointerPkgPath {
			return
		}

		targetType := pass.TypesInfo.TypeOf(argCall.Fun)
		if targetType == nil {
			return
		}

		named, ok := targetType.(*types.Named)
		if !ok || !helper.IsAzureSDKEnumType(pass, named) {
			return
		}

		if loader.IsFileChanged(pos.Filename) && !ignorer.ShouldIgnore(azbp003Name, call) {
			reporting.Reportf(pass, reporting.ReportOptions{
				Rule:          azbp003Name,
				ReportPos:     call.Pos(),
				EvidenceFile:  pos.Filename,
				EvidenceLines: []int{pos.Line},
				MatchMode:     reporting.MatchModeExactAdded,
			}, "%s: use `%s` to convert Enum type instead of explicitly type conversion.\n",
				azbp003Name, helper.FixedCode("pointer.ToEnum"))
		}
	})

	return nil, nil
}
