package AZURERMR001

import (
	"go/ast"
	"strings"

	"github.com/bflad/tfproviderlint/helper/astutils"
	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/pkgcallexprfactory"
	"golang.org/x/tools/go/analysis"
)

const Doc = `check for fmt.Errorf() using "Error" prefix

The AZURERMR001 analyzer reports when a fmt.Errorf() call contains the
beginning string "Error". This is redundent in context of terraform provider
since terraform itself already print an "[Error]" prefix at the beginning of 
error message.
`

const analyzerName = "AZURERMR001"

var fmterrorfcallexpr = pkgcallexprfactory.BuildAnalyzer("fmt", "Errorf")
var errorsnewcallexpr = pkgcallexprfactory.BuildAnalyzer("errors", "New")

var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  Doc,
	Requires: []*analysis.Analyzer{
		commentignore.Analyzer,
		fmterrorfcallexpr,
		errorsnewcallexpr,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	callExprs := pass.ResultOf[fmterrorfcallexpr].([]*ast.CallExpr)
	callExprs = append(callExprs, pass.ResultOf[errorsnewcallexpr].([]*ast.CallExpr)...)
	commentIgnore := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)

	for _, callExpr := range callExprs {
		if commentIgnore.ShouldIgnore(analyzerName, callExpr) {
			continue
		}

		formatString := astutils.ExprStringValue(callExpr.Args[0])
		if formatString == nil {
			continue
		}

		if !strings.HasPrefix(strings.ToLower(*formatString), "error ") {
			continue
		}

		pass.Reportf(callExpr.Pos(), `%s: prefer other leading words instead of "error" as error message`, analyzerName)
	}
	return nil, nil
}
