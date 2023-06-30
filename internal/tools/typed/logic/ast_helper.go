package logic

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

func identExpr(id *ast.Ident) string {
	if id != nil {
		return id.Name
	}
	return ""
}

func basicString(lit ast.Expr) string {
	if lit == nil {
		return ""
	}
	switch s := lit.(type) {
	case *ast.BasicLit:
		return strings.Trim(s.Value, "\"")
	case *ast.Ident:
		return s.Name
	}
	return ""
}

func newSelectorExpr(paths ...string) *ast.SelectorExpr {
	var res = &ast.SelectorExpr{}
	// paths := strings.Split(xpath, ".")
	if len(paths) < 2 {
		return nil
	}

	x := res
	for idx := len(paths) - 1; idx > 0; idx-- {
		item := paths[idx]
		x.Sel = ast.NewIdent(item)
		if idx == 1 {
			x.X = ast.NewIdent(paths[0])
		} else {
			sub := &ast.SelectorExpr{}
			x.X = sub
			x = sub
		}
	}
	return res
}

func (g *generator) sprintASTNode(node ast.Node, fs *token.FileSet) string {
	var buf bytes.Buffer
	cnode := &printer.CommentedNode{
		Node:     node,
		Comments: g.ast.Comments,
	}
	if err := (&printer.Config{}).Fprint(&buf, fs, cnode); err != nil {
		log.Printf("[Error] print ast with: %v", err)
	}
	return buf.String()
}

func (g *generator) sprintBody(body *ast.BlockStmt) string {
	res := g.sprintASTNode(body, g.fs)
	res = strings.TrimSuffix(strings.TrimPrefix(res, "{"), "}")
	return res
}
