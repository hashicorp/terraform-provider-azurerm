package schema

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"

	"github.com/bflad/tfproviderlint/helper/astutils"
	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

const completeSchemaDoc = `Extracts and resolves all schema fields from schema map literals.

This analyzer builds a complete view of schema maps by:
1. Parsing map[string]*pluginsdk.Schema composite literals
2. Resolving function calls to their schema definitions (cross-package, same-package, and external packages)

Resolution strategies:
- Direct literals: &pluginsdk.Schema{Type: TypeString, Required: true}
- Cross-package calls: commonschema.ResourceGroupName() → resolved from CommonAnalyzer cache
- Same-package calls: metadataSchema() → traces function definition in current package
- External package calls: network.SubnetSchema() → searches in global package registry

Output format:
- Key: token.Pos of the schema map composite literal (unique across packages in same build)
- Value: []SchemaFieldInfo with resolved schema properties for each field

Example:

    // Input code in internal/services/compute
    func (r MyResource) Arguments() map[string]*pluginsdk.Schema {
        return map[string]*pluginsdk.Schema{  // ← This map's Pos() is the cache key
            "name": {                         // ← Direct literal: resolved inline
                Type:     TypeString,
                Required: true,
            },
            "resource_group_name": commonschema.ResourceGroupName(),  // ← Cross-package: resolved from CommonAnalyzer
            "tags": tagsSchema(),                                      // ← Same-package: traces to tagsSchema() in current file
            "subnet_id": network.SubnetIdSchema(),                     // ← External package: searches in internal/services/network
        }
    }

    // Output (accessible via pass.ResultOf[CompleteSchemaAnalyzer])
    CompleteSchemaInfo.SchemaFields[mapPos] = []SchemaFieldInfo{
        {Name: "name", SchemaInfo: {Type: String, Required: true}, Position: 0},
        {Name: "resource_group_name", SchemaInfo: {Type: String, Required: true, ForceNew: true}, Position: 1},
        {Name: "tags", SchemaInfo: {...}, Position: 2},
        {Name: "subnet_id", SchemaInfo: {...}, Position: 3},
    }

Limitations:
- External package resolution depends on packages being loaded by the runner (via packages.Load)
- Does not handle dynamic schema construction (mergeSchemas, conditional schema, feature flags)
- Preserves original field order from source code

Usage:
Other analyzers can retrieve resolved schema fields without re-parsing:
    completeInfo := pass.ResultOf[CompleteSchemaAnalyzer].(*CompleteSchemaInfo)
    fields := completeInfo.SchemaFields[schemaMapNode.Pos()]
`

type CompleteSchemaInfo struct {
	// Map of schemaMapLit.Pos() -> *schema.SchemaInfo
	SchemaFields map[token.Pos][]helper.SchemaFieldInfo
}

var CompleteSchemaAnalyzer = &analysis.Analyzer{
	Name:       "completeschemainfo",
	Doc:        completeSchemaDoc,
	Run:        runComplete,
	Requires:   []*analysis.Analyzer{inspect.Analyzer, CommonAnalyzer},
	ResultType: reflect.TypeOf(&CompleteSchemaInfo{}),
}

func runComplete(pass *analysis.Pass) (interface{}, error) {
	commonSchemaInfo, ok := pass.ResultOf[CommonAnalyzer].(*CommonSchemaInfo)
	if !ok {
		return nil, nil
	}
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	completeSchemaInfo := &CompleteSchemaInfo{
		SchemaFields: make(map[token.Pos][]helper.SchemaFieldInfo),
	}

	nodeFilter := []ast.Node{(*ast.CompositeLit)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		if !helper.IsSchemaMap(comp, pass.TypesInfo) {
			return
		}

		pos := comp.Pos()
		fields := extractCompleteSchemaInfoFromMap(pass, comp, commonSchemaInfo)
		completeSchemaInfo.SchemaFields[pos] = fields
	})

	return completeSchemaInfo, nil
}

func extractCompleteSchemaInfoFromMap(pass *analysis.Pass, smap *ast.CompositeLit, commonSchemaInfo *CommonSchemaInfo) []helper.SchemaFieldInfo {
	fields := make([]helper.SchemaFieldInfo, 0, len(smap.Elts))

	for i, elt := range smap.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		fieldName := astutils.ExprStringValue(kv.Key)
		if fieldName == nil {
			continue
		}

		var resolvedSchema *schema.SchemaInfo
		switch v := kv.Value.(type) {
		case *ast.CompositeLit:
			resolvedSchema = schema.NewSchemaInfo(v, pass.TypesInfo)
		case *ast.CallExpr:
			resolvedSchema = resolveSchemaFromCall(pass, v, commonSchemaInfo)
		default:
			continue
		}

		fields = append(fields, helper.SchemaFieldInfo{
			Name:       *fieldName,
			SchemaInfo: resolvedSchema,
			Position:   i,
		})
	}

	return fields
}

// resolveSchemaFromCall resolves schema from a function call using:
// 1. CommonAnalyzer cache (commonschema.*)
// 2. Current package definitions
// 3. External packages (via global registry)
func resolveSchemaFromCall(pass *analysis.Pass, call *ast.CallExpr, commonSchemaInfo *CommonSchemaInfo) *schema.SchemaInfo {
	if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
		if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
			if obj := pass.TypesInfo.Uses[pkgIdent]; obj != nil {
				if pkgName, ok := obj.(*types.PkgName); ok {
					funcKey := pkgName.Imported().Path() + "." + selExpr.Sel.Name
					if cachedSchemaInfo, ok := commonSchemaInfo.Functions[funcKey]; ok {
						return cachedSchemaInfo
					}
				}
			}
		}
	}

	return findSchemaInCurrentPackage(pass, call)
}

// findSchemaInCurrentPackage searches current package for schema definition,
// falling back to external packages if not found.
func findSchemaInCurrentPackage(pass *analysis.Pass, call *ast.CallExpr) *schema.SchemaInfo {
	var funcObj types.Object

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		funcObj = pass.TypesInfo.Uses[fun.Sel]
	case *ast.Ident:
		funcObj = pass.TypesInfo.Uses[fun]
	default:
		return nil
	}

	if funcObj == nil {
		return nil
	}

	funcDecl := helper.FindFuncDecl(pass, funcObj)
	if funcDecl == nil {
		return findSchemaInExternalPackage(funcObj, pass.TypesInfo)
	}

	return extractSchemaFromFuncReturn(funcDecl, pass.TypesInfo)
}

// findSchemaInExternalPackage searches external packages via global registry.
func findSchemaInExternalPackage(funcObj types.Object, typesInfo *types.Info) *schema.SchemaInfo {
	if funcObj == nil || funcObj.Pkg() == nil {
		return nil
	}

	targetPkgPath := funcObj.Pkg().Path()
	funcName := funcObj.Name()

	pkg := helper.FindPackageByPath(targetPkgPath)
	if pkg == nil {
		return nil
	}

	return findSchemaInPackage(pkg, funcName, typesInfo)
}

// findSchemaInPackage searches for a function by name and extracts its schema.
func findSchemaInPackage(pkg *packages.Package, funcName string, typesInfo *types.Info) *schema.SchemaInfo {
	if pkg == nil || pkg.Syntax == nil {
		return nil
	}

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name == nil {
				continue
			}

			if funcDecl.Name.Name == funcName {
				useTypesInfo := pkg.TypesInfo
				if useTypesInfo == nil {
					useTypesInfo = typesInfo
				}
				return extractSchemaFromFuncReturn(funcDecl, useTypesInfo)
			}
		}
	}

	return nil
}
