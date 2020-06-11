package pkgcallexprfactory

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var analyzers = map[string]*analysis.Analyzer{}

func BuildAnalyzer(pkg, f string) *analysis.Analyzer {
	k := encodeAnalyzer(pkg, f)
	if analyzer, ok := analyzers[k]; ok {
		return analyzer
	}
	analyzer := &analysis.Analyzer{
		Name: fmt.Sprintf("%s%scallexpr", strings.ToLower(pkg), strings.ToLower(f)),
		Doc:  fmt.Sprintf("find %s.%s() *ast.CallExpr for later passes", pkg, f),
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
			nodeFilter := []ast.Node{
				(*ast.CallExpr)(nil),
			}
			var result []*ast.CallExpr

			inspect.Preorder(nodeFilter, func(n ast.Node) {
				callExpr := n.(*ast.CallExpr)

				switch fun := callExpr.Fun.(type) {
				case *ast.SelectorExpr:
					if fun.Sel.Name != f {
						return
					}

					switch x := fun.X.(type) {
					case *ast.Ident:
						pkgName, ok := pass.TypesInfo.ObjectOf(x).(*types.PkgName)
						if !ok {
							return
						}
						if pkgName.Imported().Path() != pkg {
							return
						}
					default:
						return
					}
				default:
					return
				}

				result = append(result, callExpr)
			})

			return result, nil
		},
		ResultType: reflect.TypeOf([]*ast.CallExpr{}),
	}
	analyzers[k] = analyzer
	return analyzer
}

func encodeAnalyzer(pkg, f string) string {
	return fmt.Sprintf("%s.%s", pkg, f)
}
