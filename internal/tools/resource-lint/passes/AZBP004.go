package passes

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const azbp004Name = "AZBP004"

const AZBP004Doc = `check for pointer.From usage

The AZBP004 analyzer reports when manual nil checks and pointer dereferencing
are used instead of pointer.From().

pointer.From returns the dereferenced value or the zero value if the pointer is nil.
Using pointer.From is more concise and handles nil cases safely.

Example violation:
  enabled := false
  if props.Enabled != nil {
      enabled = *props.Enabled
  }

Valid usage:
  enabled := pointer.From(props.Enabled)
`

var AZBP004Analyzer = &analysis.Analyzer{
	Name:     azbp004Name,
	Doc:      AZBP004Doc,
	Run:      runAZBP004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var azbp004SkipPackages = []string{"_test", "/migration", "/client", "/validate", "/test-data", "/parse", "/models"}

func runAZBP004(pass *analysis.Pass) (interface{}, error) {
	// Skip specified packages
	pkgPath := pass.Pkg.Path()
	for _, skip := range azbp004SkipPackages {
		if strings.Contains(pkgPath, skip) {
			return nil, nil
		}
	}

	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	checkZeroInitPattern(pass, inspector)

	return nil, nil
}

// checkZeroInitPattern checks for: y := <zero>; if x != nil { y = *x }
func checkZeroInitPattern(pass *analysis.Pass, inspector *inspector.Inspector) {
	nodeFilter := []ast.Node{(*ast.BlockStmt)(nil), (*ast.CaseClause)(nil)}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		var stmts []ast.Stmt

		switch node := n.(type) {
		case *ast.BlockStmt:
			stmts = node.List
		case *ast.CaseClause:
			stmts = node.Body
		default:
			return
		}

		// Look for pattern: assignment followed by if statement
		for i := 0; i < len(stmts)-1; i++ {
			// Quick check: next statement must be an if statement (fast type assertion)
			ifStmt, ok := stmts[i+1].(*ast.IfStmt)
			if !ok {
				continue
			}

			// Check if current statement is a zero-value assignment
			assignStmt, varName := isZeroValueAssignment(stmts[i], pass)
			if assignStmt == nil {
				continue
			}

			// Check the if statement matches the pattern
			if !isMatchingNilCheckAssignment(ifStmt, varName) {
				continue
			}

			pos := pass.Fset.Position(assignStmt.Pos())
			if loader.ShouldReport(pos.Filename, pos.Line) {
				pass.Reportf(assignStmt.Pos(),
					"%s: can simplify with `%s` since variable is initialized to zero value\n",
					azbp004Name, helper.FixedCode("pointer.From()"))
			}
		}
	})
}

// isZeroValueAssignment checks if a statement is: varName := <zero-value>
// Returns the assignment statement and variable name if it matches
func isZeroValueAssignment(stmt ast.Stmt, pass *analysis.Pass) (*ast.AssignStmt, string) {
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return nil, ""
	}

	// Must be := (short variable declaration)
	if assignStmt.Tok != token.DEFINE {
		return nil, ""
	}

	// Must have exactly one LHS and one RHS
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return nil, ""
	}

	// LHS must be an identifier
	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, ""
	}

	// Check if RHS is a zero value
	if !isZeroValue(assignStmt.Rhs[0], pass) {
		return nil, ""
	}

	return assignStmt, lhsIdent.Name
}

// isZeroValue checks if an expression is a zero value (false, 0, "", nil, etc.)
func isZeroValue(expr ast.Expr, pass *analysis.Pass) bool {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok {
		return false
	}

	// Check if it's a constant with zero value
	if tv.Value != nil {
		switch tv.Value.Kind() {
		case constant.Bool:
			return !constant.BoolVal(tv.Value)
		case constant.String:
			return constant.StringVal(tv.Value) == ""
		case constant.Int:
			return constant.Sign(tv.Value) == 0
		case constant.Float:
			return constant.Sign(tv.Value) == 0
		}
	}

	// Check for nil
	if isNilIdent(expr) {
		return true
	}

	return false
}

// isMatchingNilCheckAssignment checks if ifStmt matches: if ptrExpr != nil { varName = *ptrExpr }
func isMatchingNilCheckAssignment(ifStmt *ast.IfStmt, varName string) bool {
	// No else branch
	if ifStmt.Else != nil {
		return false
	}

	// Condition must be "expr != nil"
	binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return false
	}

	// Check right side is nil
	nilIdent, ok := binExpr.Y.(*ast.Ident)
	if !ok || nilIdent.Name != "nil" {
		return false
	}

	checkedExpr := binExpr.X

	// Body must have exactly one statement
	if len(ifStmt.Body.List) != 1 {
		return false
	}

	stmt := ifStmt.Body.List[0]

	// Must be an assignment (not a declaration)
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok || len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return false
	}

	// Must be = (ASSIGN), not := (DEFINE)
	if assignStmt.Tok != token.ASSIGN {
		return false
	}

	// LHS must be the variable we're tracking
	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok || lhsIdent.Name != varName {
		return false
	}

	// RHS must be *expr (StarExpr in AST, not UnaryExpr)
	starExpr, ok := assignStmt.Rhs[0].(*ast.StarExpr)
	if !ok {
		return false
	}

	// The dereferenced expression must match the nil-checked expression
	return astExprEqual(checkedExpr, starExpr.X)
}

// isNilIdent checks if an expression is the nil identifier
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// astExprEqual checks if two AST expressions are structurally equal
func astExprEqual(a, b ast.Expr) bool {
	// Use types.ExprString for simple comparison
	// This handles most cases including field selectors like input.Name
	return types.ExprString(a) == types.ExprString(b)
}
