package AZURERMR001

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strconv"
	"strings"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/pkgcallexprfactory"
	"golang.org/x/tools/go/analysis"
)

const Doc = `check for error message using "Error" prefix

The AZURERMR001 analyzer reports cases of  a error message contains the
beginning string "Error". This is redundant in context of terraform provider
since terraform itself already print an "[Error]" prefix at the beginning of 
error message.
`
const errorWord = "error " // suffix space is word boundary

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

		// errors.New() and fmt.Errorf() at least has one parameter,
		// hence no need to check boundary.
		firstArg, ok := callExpr.Args[0].(*ast.BasicLit)
		if !ok {
			continue
		}
		if firstArg.Kind != token.STRING {
			continue
		}
		firstArgValue, _ := strconv.Unquote(firstArg.Value) // can assume well-formed Go

		if !strings.HasPrefix(strings.ToLower(firstArgValue), errorWord) {
			continue
		}

		// suggested fix
		var callExprBuf bytes.Buffer
		firstArg.Value = string(firstArg.Value[0]) + firstArgValue[len(errorWord):] + string(firstArg.Value[len(firstArg.Value)-1])

		if err := format.Node(&callExprBuf, pass.Fset, callExpr); err != nil {
			return nil, fmt.Errorf("error formatting new expression: %s", err)
		}

		pass.Report(analysis.Diagnostic{
			Pos:     callExpr.Pos(),
			End:     callExpr.End(),
			Message: fmt.Sprintf(`%s: prefer other leading words instead of "error" as error message`, analyzerName),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Remove",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     callExpr.Pos(),
							End:     callExpr.End(),
							NewText: callExprBuf.Bytes(),
						},
					},
				},
			},
		})
	}
	return nil, nil
}
