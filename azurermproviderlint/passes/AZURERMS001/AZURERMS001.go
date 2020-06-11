package AZURERMS001

import (
	"bytes"
	"fmt"
	"github.com/bflad/tfproviderlint/helper/astutils"
	"go/ast"
	"go/format"
	"golang.org/x/tools/go/analysis"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/bflad/tfproviderlint/passes/helper/schema/schemainfo"
)

const Doc = `check for Schema contains case-insensitive validation missing case diff suppression

The AZURERMS001 analyzer reports cases of schema that uses case-insensitive
validation missing case diff suppression.`

const analyzerName = "AZURERMS001"

var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  Doc,
	Requires: []*analysis.Analyzer{
		schemainfo.Analyzer,
		commentignore.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ignorer := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	schemaInfos := pass.ResultOf[schemainfo.Analyzer].([]*schema.SchemaInfo)
	for _, schemaInfo := range schemaInfos {
		if ignorer.ShouldIgnore(analyzerName, schemaInfo.AstCompositeLit) {
			continue
		}

		// Ignore schemas which haven't declared `ValidateFunc: validation.StringInSlice($_, true)`
		if !schemaInfo.DeclaresField(schema.SchemaFieldValidateFunc) {
			continue
		}
		switch validateFuncExpr := schemaInfo.Fields[schema.SchemaFieldValidateFunc].Value.(type) {
		case *ast.CallExpr:
			if !astutils.IsPackageFunctionFieldListType(validateFuncExpr.Fun, pass.TypesInfo, "validation", "StringInSlice") {
				continue
			}
			if len(validateFuncExpr.Args) != 2 {
				continue
			}
			arg, ok := validateFuncExpr.Args[1].(*ast.Ident)
			if !ok {
				continue
			}
			if arg.Name != "true" {
				continue
			}
		default:
			continue
		}

		// Check wether this schema has defined `DiffSuppressFunc: suppress.CaseDifference`
		if diffsuppressFuncField := schemaInfo.Fields[schema.SchemaFieldDiffSuppressFunc]; diffsuppressFuncField != nil {
			if diffsuppressFuncExpr, ok := diffsuppressFuncField.Value.(*ast.SelectorExpr); ok {
				if astutils.IsPackageFunctionFieldListType(diffsuppressFuncExpr, pass.TypesInfo, "suppress", "CaseDifference") {
					continue
				}
			}
		}

		// suggested fix
		var fix *analysis.SuggestedFix
		if diffsuppressFuncField := schemaInfo.Fields[schema.SchemaFieldDiffSuppressFunc]; diffsuppressFuncField != nil {
			// Replace
			e := &ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "DiffSuppressFunc",
				},
				Value: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "suppress"},
					Sel: &ast.Ident{Name: "CaseDifference"},
				},
			}
			var buf bytes.Buffer
			if err := format.Node(&buf, pass.Fset, e); err != nil {
				return nil, fmt.Errorf("error formatting new expression: %s", err)
			}
			fix = &analysis.SuggestedFix{
				Message: "Replace",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     diffsuppressFuncField.Pos(),
						End:     diffsuppressFuncField.End(),
						NewText: buf.Bytes(),
					},
				},
			}
		} else {
			// Add
			e := &ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "DiffSuppressFunc",
				},
				Value: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "suppress",
					},
					Sel: &ast.Ident{
						Name: "CaseDifference",
					},
				},
			}
			var buf bytes.Buffer
			if err := format.Node(&buf, pass.Fset, e); err != nil {
				return nil, fmt.Errorf("error formatting new expression: %s", err)
			}
			pos := schemaInfo.AstCompositeLit.Rbrace
			fix = &analysis.SuggestedFix{
				Message: "Add",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     pos,
						NewText: buf.Bytes(),
					},
				},
			}
		}

		pass.Report(analysis.Diagnostic{
			Pos:            schemaInfo.AstCompositeLit.Pos(),
			End:            schemaInfo.AstCompositeLit.End(),
			Message:        fmt.Sprintf("%s: prefer adding `DiffSuppressFunc: suppress.CaseDifference` when ignoring case during validation", analyzerName),
			SuggestedFixes: []analysis.SuggestedFix{*fix},
		})
	}

	return nil, nil
}
