package logic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// from an untyped resource, generate source code of a typed resource
// logic design
type generator struct {
	source  *untypedResource
	builder *strings.Builder
	fs      *token.FileSet
	ast     *ast.File
	schemas map[string]ast.Expr

	meta *MetaInfo
}

func newGenerator(untyped *untypedResource) (*generator, error) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, untyped.path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	ast.SortImports(fs, file)
	return &generator{
		source:  untyped,
		builder: &strings.Builder{},
		fs:      fs,
		ast:     file,
		schemas: map[string]ast.Expr{},
	}, nil
}

func (g *generator) buildMeta() *MetaInfo {
	meta := g.metaInfo()
	// function body transform to after g.meta assigned
	for _, fn := range meta.Functions {
		g.transformFn(fn)
		fn.BodyContent = g.sprintBody(fn.Ast.Body)
	}
	for _, fn := range []*Function{meta.CreateFunc, meta.UpdateFunc, meta.ReadFunc, meta.DeleteFunc} {
		if fn == nil {
			continue // Update can be nil
		}
		g.transformFn(fn)
		fn.BodyContent = g.sprintBody(fn.Ast.Body)
	}
	return meta
}

func (g *generator) metaInfo() *MetaInfo {
	if g.meta == nil {
		meta := g.collectASTInfo()
		g.buildModels(meta)
		g.meta = meta

	}
	return g.meta
}

func (g *generator) collectASTInfo() (res *MetaInfo) {
	res = &MetaInfo{
		Name:         camelCase(g.source.name) + "Resource",
		Package:      identExpr(g.ast.Name),
		ResourceType: g.source.name,
		Recv:         g.source.name[:1],
	}
	// collect read/write function name
	fnName := func(fn interface{}) string {
		if fn == nil {
			return ""
		}
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		if idx := strings.LastIndexByte(name, '.'); idx > 0 {
			return name[idx+1:]
		}
		return name
	}
	r := g.source.resource

	createName := fnName(r.Create) //nolint:staticcheck
	readName := fnName(r.Read)     //nolint:staticcheck
	updateName := fnName(r.Update) //nolint:staticcheck
	deleteName := fnName(r.Delete) //nolint:staticcheck

	// collect function asts
	ast.Inspect(g.ast, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ImportSpec:
			res.Imports = append(res.Imports, n)
		case *ast.FuncDecl:
			f := &Function{
				Name:        n.Name.String(),
				Timout:      0,
				Ast:         n,
				BodyContent: g.sprintBody(n.Body),
			}
			switch f.Name {
			case createName:
				res.CreateFunc = f
				f.Timout = pointer.From(r.Timeouts.Create)
			case readName:
				res.ReadFunc = f
				f.Timout = pointer.From(r.Timeouts.Read)
				// we usually use `model` variable for `sdkResponse.Model`, so rename the model name
				f.varModel = "stateModel"
			case updateName:
				res.HasUpdate = true
				res.UpdateFunc = f
				f.Timout = pointer.From(r.Timeouts.Update)
			case deleteName:
				res.DeleteFunc = f
				f.Timout = pointer.From(r.Timeouts.Delete)
			default:
				if f.Name != g.source.funcName {
					res.Functions = append(res.Functions, f)
				}
			}

			if fnRes := n.Type.Results; fnRes.NumFields() == 1 {
				// resource define function
				if resType := types.ExprString(fnRes.List[0].Type); resType == "*pluginsdk.Resource" {
					// inspect resource in top level pluginsdk.Resource
					ast.Inspect(n.Body, func(node ast.Node) bool {
						if c, ok := node.(*ast.CompositeLit); ok {
							if cType := types.ExprString(c.Type); cType == "pluginsdk.Resource" {
								// compositeLit should contain Create
								if g.collectResourceDefine(c, res) {
									return false
								}
							}
						}
						return true
					})
					// inspect schema
					ast.Inspect(n.Body, func(node ast.Node) bool {
						if c, ok := node.(*ast.CompositeLit); ok {
							// top level schema definition
							if cType := types.ExprString(c.Type); cType == "map[string]*pluginsdk.Schema" {
								g.collectSchemaFromAST(c)
								return false
							}
						}
						return true
					})
				}
			}
			return false
		}
		return true
	})

	g.splitSchema(res)
	return res
}

func (g *generator) collectResourceDefine(c *ast.CompositeLit, meta *MetaInfo) (isResource bool) {
	for _, ele := range c.Elts {
		if st, ok := ele.(*ast.KeyValueExpr); ok {
			key := basicString(st.Key)
			switch key {
			case "Create":
				isResource = true
			case "Importer":
				meta.IDValidator = g.getImporter(st.Value)
			case "SchemaVersion":
				meta.SchemaVersion = basicString(st.Value)
			case "StateUpgraders":
				meta.StateUpgrade = g.getStateUpgrade(st.Value)
			}
		}
	}
	return isResource
}

func (g *generator) collectSchemaFromAST(c *ast.CompositeLit) {
	for _, ele := range c.Elts {
		if st, ok := ele.(*ast.KeyValueExpr); ok {
			key := basicString(st.Key)
			if _, ok := g.schemas[key]; !ok {
				g.schemas[key] = st.Value
			}
		}
	}
}

// call after collected all schemas
func (g *generator) splitSchema(meta *MetaInfo) {
	// range over top-level schema
	for k, sch := range g.source.resource.Schema {
		fromAst, ok := g.schemas[k]
		if !ok {
			log.Printf("no schema `%s` found from ast, skip", k)
			continue
		}
		content := g.sprintASTNode(fromAst, g.fs)
		if sch.Required || sch.Optional {
			meta.Arguments = append(meta.Arguments, NewSchemaFieldFromAST(k, content))
		} else {
			meta.Attributes = append(meta.Attributes, NewSchemaFieldFromAST(k, content))
		}
	}
	// sort by name
	sort.Slice(meta.Arguments, func(i, j int) bool {
		return meta.Arguments[i].Name < meta.Arguments[j].Name
	})
	sort.Slice(meta.Attributes, func(i, j int) bool {
		return meta.Attributes[i].Name < meta.Attributes[j].Name
	})
}

func (g *generator) getImporter(expr ast.Expr) (idValidator string) {
	ast.Inspect(expr, func(node ast.Node) bool {
		// try to expand the first function call of id validate
		// if a function call arg[0] is id, them use it
		if n, ok := node.(*ast.CallExpr); ok {
			func() {
				defer func() {
					_ = recover()
				}()
				if n.Args[0].(*ast.Ident).Name == "id" {
					idValidator = types.ExprString(n.Fun)
				}
			}()
			if idValidator != "" {
				idValidator = strings.ReplaceAll(idValidator, "Parse", "Validate")
				return false
			}
		}
		return true
	})
	return
}

func (g *generator) getStateUpgrade(expr ast.Expr) (res string) {
	ast.Inspect(expr, func(node ast.Node) bool {
		if l, ok := node.(*ast.CompositeLit); ok {
			if typ := types.ExprString(l.Type); typ == "map[int]pluginsdk.StateUpgrade" {
				res = types.ExprString(l)
			}
		}
		return res == ""
	})
	return res
}

func (g *generator) buildModels(meta *MetaInfo) {
	name := strings.TrimPrefix(g.source.name, "azurerm_")
	meta.ModelName = modelName(name)
	g.buildSchema(name, g.source.resource.Schema, meta)
}

func (g *generator) buildSchema(key string, schema map[string]*pluginsdk.Schema, meta *MetaInfo) {
	m := Model{
		Name: modelName(key),
	}
	for name, sch := range schema {
		// identity field
		if f := g.specialSchema(name); f != nil {
			m.Fields = append(m.Fields, *f)
			continue
		}
		m.Fields = append(m.Fields, ModelField{
			Name: camelCase(name),
			Type: valueType(sch, name),
			Tag:  name,
		})
		if res, ok := sch.Elem.(*pluginsdk.Resource); ok {
			g.buildSchema(name, res.Schema, meta)
		}
	}
	meta.Models = append(meta.Models, m)
}

func (g *generator) specialSchema(key string) *ModelField {
	expr, ok := g.schemas[key]
	if !ok {
		return nil
	}
	var f = &ModelField{
		Name: camelCase(key),
		Tag:  key,
	}
	if key == "identity" {
		if call, ok := expr.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				name := sel.Sel.Name
				switch {
				case strings.Contains(name, "SystemAssignedUserAssignedIdentity"):
					f.Type = "[]identity.ModelSystemAssignedUserAssigned"
				case strings.Contains(name, "SystemOrUserAssignedIdentity"):
					f.Type = "[]identity.ModelSystemAssignedUserAssigned"
				case strings.Contains(name, "SystemAssignedIdentity"):
					f.Type = "[]identity.ModelSystemAssigned"
				case strings.Contains(name, "UserAssignedIdentity"):
					f.Type = "[]identity.ModelUserAssigned"
				default:
					f.Type = "[]identity.ModelSystemAssignedUserAssigned"
				}
				return f
			}
		}
	}
	return nil
}
