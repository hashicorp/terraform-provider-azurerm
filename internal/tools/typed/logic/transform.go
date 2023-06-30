package logic

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
)

// shared between tranforms
type funcInfo struct {
	dname    string // parameter name of pluginsdk.ResourceData, d in most of the time
	varModel string // model variable name
	fn       *ast.FuncDecl
}

// used to transform function body
type transformer interface {
	transform(info *funcInfo, node ast.Node, g *generator)
}

type fnAST struct {
	file    *ast.File
	fs      *token.FileSet
	content string
}

var (
	// by stage 0, 1, 2...
	transformers = [][]transformer{
		{
			metaTransform{},
			dGetTransform{},
			dSetTransform{},
			importTransform{},
			stmtTransform{},
			fnNameReplace{},
		},
		{
			// inlineFuncTransform should after dGet, dSet all done
			inlineFuncTransform{},
			literalMapToModelTransform{},
		},
	}
)

func (g *generator) transformFn(fun *Function) {
	if fun.Ast == nil {
		return
	}
	fn := fun.Ast
	var fnInfo = &funcInfo{
		varModel: fun.VarModel(),
		fn:       fn,
	}
	if params := fn.Type.Params.List; len(params) > 0 {
		if _, ok := params[0].Type.(*ast.StarExpr); ok {
			if typ := types.ExprString(params[0].Type); typ == "*pluginsdk.ResourceData" {
				fnInfo.dname = params[0].Names[0].Name
			}
		}
	}

	for _, tt := range transformers {
		ast.Inspect(fn, func(node ast.Node) bool {
			for _, tf := range tt {
				tf.transform(fnInfo, node, g)
			}
			return true
		})
	}
}

type stmtTransform struct{}

func (s stmtTransform) transform(fnInfo *funcInfo, node ast.Node, g *generator) {
	body, ok := node.(*ast.BlockStmt)
	if !ok {
		return
	}

	var list []ast.Stmt
	for idx := 0; idx < len(body.List); idx++ {
		stmt := body.List[idx]
		// ctx, cancel := timeouts.ForXXX()
		if ass, ok := stmt.(*ast.AssignStmt); ok {
			if call, ok := ass.Rhs[0].(*ast.CallExpr); ok {
				if str := types.ExprString(call.Fun); strings.HasPrefix(str, "timeouts.For") {
					// remove this line, or success line `defer cancel()`
					if idx+1 < len(body.List) {
						if _, ok := body.List[idx+1].(*ast.DeferStmt); ok {
							idx++
						}
					}
					continue
				}
			}
		}

		// d.SetId(id.ID()) ID change to new call => meta.SetID(id)
		if expr, ok := stmt.(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && types.ExprString(fun) == fnInfo.dname+".SetId" {
					if lit, ok := call.Args[0].(*ast.BasicLit); ok && basicString(lit) == "" {
						// modify to MarkAsGone
						newStmt := &ast.ReturnStmt{
							Return: 0,
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun:  newSelectorExpr("meta", "MarkAsGone"),
									Args: []ast.Expr{ast.NewIdent("id")},
								},
							},
						}
						// if success with a return nil, remove it
						list = append(list, newStmt)
						break // break out after MarkAsGone
					} else {
						fun.X = ast.NewIdent("meta")
						fun.Sel = ast.NewIdent("SetID")
						call.Args = []ast.Expr{ast.NewIdent("id")}
						list = append(list, stmt)
					}
					continue
				}
			}
		}

		// return xxxRead(d, meta), change to return nil
		if ret, ok := stmt.(*ast.ReturnStmt); ok && len(ret.Results) > 0 {
			if call, ok := ret.Results[0].(*ast.CallExpr); ok {
				if fnID, ok := call.Fun.(*ast.Ident); ok && strings.Contains(fnID.Name, "Read") {
					if basicString(call.Args[0]) == fnInfo.dname {
						ret.Results = []ast.Expr{ast.NewIdent("nil")}
						list = append(list, stmt)
						continue
					}
				}
			}
		}

		list = append(list, stmt)
	}
	body.List = list
}

type dSetTransform struct{}

func (d dSetTransform) transform(info *funcInfo, node ast.Node, g *generator) {
	if blocks, ok := node.(*ast.BlockStmt); ok {
		// range over blocks
		for idx, stmt := range blocks.List {
			// d.Set("") => metadata.MarkAsGone
			var call *ast.CallExpr
			switch st := stmt.(type) {
			case *ast.IfStmt:
				// if err := d.Set()
				safeRun(func() {
					call = st.Init.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr)
				})
			case *ast.ExprStmt:
				// d.Set()
				if c, ok := st.X.(*ast.CallExpr); ok {
					call = c
				}
			}
			if call != nil {
				var setCall *ast.CallExpr
				safeRun(func() {
					fun := call.Fun.(*ast.SelectorExpr)
					if fun.X.(*ast.Ident).Name == info.dname {
						if fun.Sel.Name == "Set" {
							setCall = call
						}
					}
				})
				if setCall != nil {
					key := basicString(setCall.Args[0])
					fieldName := g.meta.findFieldByTag("", key).Name
					arg := setCall.Args[1]
					// todo: check if arg is a pointer
					if sel, ok := arg.(*ast.SelectorExpr); ok && basicString(sel.X) != "id" {
						arg = &ast.CallExpr{
							Fun:  newSelectorExpr("pointer", "From"),
							Args: []ast.Expr{sel},
						}
					}
					newStmt := &ast.AssignStmt{
						Lhs: []ast.Expr{
							newSelectorExpr(info.varModel, fieldName),
						},
						// most of the time, we need a pointer.From() wrap
						Rhs: []ast.Expr{arg},
						Tok: token.ASSIGN,
					}
					blocks.List[idx] = newStmt
				}
			}

		}
	}
}

type metaTransform struct{}

func (m metaTransform) transform(info *funcInfo, node ast.Node, _ *generator) {
	if sel, ok := node.(*ast.SelectorExpr); ok {
		if assert, ok := sel.X.(*ast.TypeAssertExpr); ok {
			if types.ExprString(assert.Type) == "*clients.Client" {
				meta := basicString(assert.X)
				sel.X = &ast.SelectorExpr{
					X:   ast.NewIdent(meta),
					Sel: ast.NewIdent("Client"),
				}
			}
		}
	}
}

type dGetTransform struct{}

// in a := d.Get("key").(string) or fn(d.Get("key").(string)) or d.Id(),
// or `a, err := d.GetOk("key");`
func (d dGetTransform) transform(info *funcInfo, node ast.Node, g *generator) {
	// get the key of d.Get
	trans := func(r ast.Expr) ast.Expr {
		var call *ast.CallExpr
		switch item := r.(type) {
		case *ast.TypeAssertExpr:
			if c, ok := item.X.(*ast.CallExpr); ok {
				call = c
			}
		case *ast.CallExpr:
			call = item
		}
		if call == nil {
			return nil
		}

		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if basicString(sel.X) == info.dname {
				if selName := sel.Sel.Name; selName == "Get" {
					key := basicString(call.Args[0])
					return &ast.SelectorExpr{
						X:   ast.NewIdent(info.varModel),
						Sel: ast.NewIdent(g.metaInfo().findFieldByTag("", key).Name),
					}
				} else if selName == "GetOk" {
					// GetOk should be used in an assignment
					sel.X = &ast.SelectorExpr{
						X:   ast.NewIdent("meta"),
						Sel: ast.NewIdent("ResourceData"),
					}
					return nil // update in-place, no need to return the node
				}
			}
		}
		return nil
	}

	switch t := node.(type) {
	case *ast.SelectorExpr:
		// d.Id => meta.ResourceData.Id
		if str := types.ExprString(t); str == fmt.Sprintf("%s.Id", info.dname) {
			t.X = &ast.SelectorExpr{
				X:   ast.NewIdent("meta"),
				Sel: ast.NewIdent("ResourceData"),
			}
		}
	case *ast.AssignStmt:
		for idx, r := range t.Rhs {
			if newExpr := trans(r); newExpr != nil {
				t.Rhs[idx] = newExpr
			}
		}
	case *ast.CallExpr:
		// d.Get as an arg: fn(d.Get("").(string))
		for idx, arg := range t.Args {
			if newExpre := trans(arg); newExpre != nil {
				t.Args[idx] = newExpre
			}
		}
	}
}

type importTransform struct{}

func (i importTransform) transform(info *funcInfo, node ast.Node, g *generator) {
	switch n := node.(type) {
	case *ast.CallExpr:
		if fn, ok := n.Fun.(*ast.SelectorExpr); ok {
			if types.ExprString(fn) == "tf.ImportAsExistsError" {
				// imports.ForSchema => meta.ForSchema
				fn.X = ast.NewIdent("meta")
				fn.Sel = ast.NewIdent("ResourceRequiresImport")
				n.Args = []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(g.meta.Recv),
							Sel: ast.NewIdent("ResourceType"),
						},
					},
					ast.NewIdent("id"),
				}
			}
		}
	}
}

type fnNameReplace struct{}

var (
	callNameReplace = map[string]string{
		"identity.ExpandSystemAndUserAssignedMap":        "identity.ExpandSystemAndUserAssignedMapFromModel",
		"identity.FlattenSystemAndUserAssignedMap":       "identity.FlattenSystemAndUserAssignedMapToModel",
		"identity.ExpandSystemOrUserAssignedMap":         "identity.ExpandSystemOrUserAssignedMapFromModel",
		"identity.FlattenSystemOrUserAssignedMap":        "identity.FlattenSystemOrUserAssignedMapToModel",
		"identity.ExpandLegacySystemAndUserAssignedMap":  "identity.ExpandLegacySystemAndUserAssignedMapFromModel",
		"identity.FlattenLegacySystemAndUserAssignedMap": "identity.FlattenLegacySystemAndUserAssignedMapToModel",
		"identity.ExpandSystemAssigned":                  "identity.ExpandSystemAssignedFromModel",
		"identity.FlattenSystemAssigned":                 "identity.FlattenSystemAssignedToModel",
		"identity.ExpandUserAssignedMap":                 "identity.ExpandUserAssignedMapFromModel",
		"identity.FlattenUserAssignedMap":                "identity.FlattenUserAssignedMapToModel",
	}
)

func (f fnNameReplace) transform(info *funcInfo, node ast.Node, g *generator) {
	if call, ok := node.(*ast.CallExpr); ok {
		if fn, ok := call.Fun.(*ast.SelectorExpr); ok {
			funStr := types.ExprString(fn)
			if newName, ok := callNameReplace[funStr]; ok {
				call.Fun = newSelectorExpr(strings.Split(newName, ".")...)
			}
			if funStr == info.dname+".IsNewResource" {
				call.Fun = newSelectorExpr("meta", "ResourceData", "IsNewResource")
			}
		}
	}
}

// this transform should place after dGetTransform
type inlineFuncTransform struct{}

// transform `x := flattenXXX(prop.Model.Field)` => `x := r.flattenXXX(prop.Model.Field)`
// and `x := expandXXX(d.Get("key").([]interface{}))` => `x := r.expandXXX(model.Key)`
// also refactor flatten/expand functin signature:
//
//	flatten: change return value type
//	expand: change argument type
func (i inlineFuncTransform) transform(info *funcInfo, node ast.Node, g *generator) {
	// if call function defined in source code file, then transform to call with receiver of current model resource
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return
	}
	if len(assign.Rhs) != 1 {
		return
	}
	if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
		if id, ok := call.Fun.(*ast.Ident); ok {
			for _, fn := range g.meta.Functions {
				if fn.Name == id.Name {
					call.Fun = newSelectorExpr(g.meta.Recv, fn.Name)
					// args transform
					for idx, arg := range call.Args {
						switch a := arg.(type) {
						case *ast.Ident:
							if identExpr(a) == info.dname {
								call.Args[idx] = &ast.UnaryExpr{
									X:  ast.NewIdent(info.varModel),
									Op: token.AND,
								}
								// change function argument type as pointer
								transformFunction(fn, idx, "model", "*"+g.meta.ModelName)
							}
						case *ast.SelectorExpr:
							if xStr := basicString(a.X); xStr == info.varModel {
								typ := modelName(a.Sel.Name) // should be the same as model struct type
								transformFunction(fn, idx, "", typ)
							}
						}
					}
					// return value, change flatten function return type
					// if assign to model.Key: a flatten function call
					if lhs, ok := assign.Lhs[0].(*ast.SelectorExpr); ok {
						if x, ok := lhs.X.(*ast.Ident); ok && x.Name == info.varModel {
							typ := g.meta.topField(lhs.Sel.Name)
							if typ == nil {
								log.Printf("[Error] no top field of: %s in %s", lhs.Sel.Name, g.meta.ModelName)
							} else {
								ret := fn.Ast.Type.Results.List[0]
								if strings.HasPrefix(typ.Type, "[]") {
									ret.Type = &ast.ArrayType{
										Elt: ast.NewIdent(typ.Type[2:]),
									}
								} else {
									ret.Type = ast.NewIdent(typ.Type)
								}
							}
						}
					}
				}
			}
		}
	}
}

func transformFunction(fn *Function, idx int, toName, toTyp string) {
	if !fn.transformed {
		fn.transformed = true
		arg0 := fn.Ast.Type.Params.List[idx]
		if arr, ok := arg0.Type.(*ast.ArrayType); ok {
			arr.Elt = ast.NewIdent(toTyp)
		} else {
			arg0.Type = ast.NewIdent(toTyp)
		}
		if toName != "" {
			arg0.Names[0] = ast.NewIdent(toName)
		}
		// trans to
	}
}

type literalMapToModelTransform struct{}

// as for return map[string]interface{} or return []interface{}{map[string]interface{}{}}
func (l literalMapToModelTransform) transform(info *funcInfo, node ast.Node, g *generator) {
	if comp, ok := node.(*ast.CompositeLit); ok {
		typ := types.ExprString(comp.Type)
		if typ == "map[string]interface{}" {
			// extract keys
			m := map[string]ast.Expr{}
			var keys []string
			for _, el := range comp.Elts {
				if kv, ok := el.(*ast.KeyValueExpr); ok {
					if key, ok := kv.Key.(*ast.BasicLit); ok {
						m[basicString(key)] = kv.Value
						keys = append(keys, basicString(key))
					}
				} else {
					return
				}
			}

			if model := g.meta.findModelByKeys(keys); model != nil {
				comp.Type = ast.NewIdent(model.Name)
				var els []ast.Expr
				for _, f := range model.Fields {
					if v, ok := m[f.Tag]; ok {
						els = append(els, &ast.KeyValueExpr{
							Key:   ast.NewIdent(f.Name),
							Value: v,
						})
					}
				}
				comp.Elts = els
			}
		}
	}
}

var (
	_ = []transformer{
		metaTransform{},
		dGetTransform{},
		dSetTransform{},
		importTransform{},
		fnNameReplace{},
	}
)
