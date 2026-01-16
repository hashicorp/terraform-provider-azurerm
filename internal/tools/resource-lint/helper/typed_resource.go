package helper

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	PackagePathSDK = "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// TypedResourceInfo represents gathered information about a typed Terraform resource
type TypedResourceInfo struct {
	ResourceTypeName     string
	ModelName            string
	ModelStruct          *ast.StructType
	ArgumentsFunc        *ast.FuncDecl
	ArgumentsProperties  []SchemaFieldInfo // Parsed schema fields from Arguments()
	AttributesFunc       *ast.FuncDecl
	CreateFunc           *ast.FuncDecl
	ReadFunc             *ast.FuncDecl
	UpdateFunc           *ast.FuncDecl
	UpdateFuncBody       *ast.BlockStmt
	DeleteFunc           *ast.FuncDecl
	TypesInfo            *types.Info
	ModelFieldToTFSchema map[string]string // model struct field name -> tfschema tag name
}

// NewTypedResourceInfo creates a TypedResourceInfo by parsing a typed resource from file
func NewTypedResourceInfo(resourceTypeName string, file *ast.File, info *types.Info) *TypedResourceInfo {
	result := &TypedResourceInfo{
		ResourceTypeName:     resourceTypeName,
		TypesInfo:            info,
		ModelFieldToTFSchema: make(map[string]string),
	}

	// Single pass: collect all information from file.Decls
	structs := make(map[string]*ast.StructType)
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// Collect struct definitions
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structs[typeSpec.Name.Name] = structType
					}
				}
			}

		case *ast.FuncDecl:
			if d.Recv == nil || len(d.Recv.List) == 0 {
				continue
			}

			recvType := GetReceiverTypeName(d.Recv.List[0].Type)
			if recvType != resourceTypeName {
				continue
			}

			// Collect methods by name
			switch d.Name.Name {
			case "ModelObject":
				// Extract model type name from: return &ModelName{}
				if d.Body != nil {
					ast.Inspect(d.Body, func(n ast.Node) bool {
						ret, ok := n.(*ast.ReturnStmt)
						if !ok || len(ret.Results) == 0 {
							return true
						}
						// Match &ModelName{} pattern
						if unaryExpr, ok := ret.Results[0].(*ast.UnaryExpr); ok {
							if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
								if ident, ok := compLit.Type.(*ast.Ident); ok {
									result.ModelName = ident.Name
									return false
								}
							}
						}
						return true
					})
				}

			case "Arguments":
				result.ArgumentsFunc = d
			case "Attributes":
				result.AttributesFunc = d
			case "Create":
				result.CreateFunc = d
			case "Read":
				result.ReadFunc = d
			case "Update":
				result.UpdateFunc = d
			case "Delete":
				result.DeleteFunc = d
			}
		}
	}

	// Resolve model struct from collected structs
	if result.ModelName != "" {
		if modelStruct, ok := structs[result.ModelName]; ok {
			result.ModelStruct = modelStruct
			result.parseModelStruct(modelStruct)
		}
	}

	return result
}

// parsing model struct
func (info *TypedResourceInfo) parseModelStruct(modelStruct *ast.StructType) {
	for _, field := range modelStruct.Fields.List {
		if field.Tag == nil {
			continue
		}

		tagValue := strings.Trim(field.Tag.Value, "`")
		if !strings.Contains(tagValue, "tfschema:") {
			continue
		}

		// Extract tfschema tag value
		parts := strings.Split(tagValue, `"`)
		if len(parts) >= 2 {
			tfschemaName := parts[1]
			if len(field.Names) > 0 {
				// Map: struct field name -> tfschema name
				info.ModelFieldToTFSchema[field.Names[0].Name] = tfschemaName
			}
		}
	}
}

// GetReceiverTypeName extracts the type name from a method receiver
func GetReceiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// extractFuncBodyFromCompositeLit extracts Func body from sdk.ResourceFunc{Func: func(){...}}
// Validates that the composite literal is of type sdk.ResourceFunc
func extractFuncBodyFromCompositeLit(compLit *ast.CompositeLit, typesInfo *types.Info) *ast.BlockStmt {
	typ := typesInfo.TypeOf(compLit)
	// Check type name ends with "ResourceFunc" to support both real SDK and test mock
	if typ != nil {
		if named, ok := typ.(*types.Named); !ok || !strings.HasSuffix(named.Obj().Name(), "ResourceFunc") {
			return nil
		}
	}

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == "Func" {
			if funcLit, ok := kv.Value.(*ast.FuncLit); ok {
				return funcLit.Body
			}
		}
	}
	return nil
}

// GetFuncBody extracts the actual function body from sdk.ResourceFunc
// Handles three cases:
// 1. Direct: return sdk.ResourceFunc{Func: func(){...}}
// 2. Base method: return r.base.someFunc()
// 3. Helper function: return someHelper()
func GetFuncBody(pass *analysis.Pass, funcDecl *ast.FuncDecl) *ast.BlockStmt {
	if funcDecl == nil || funcDecl.Body == nil {
		return nil
	}

	var funcBody *ast.BlockStmt
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			return true
		}

		// Case 1: Direct sdk.ResourceFunc{Func: func(){...}}
		if compLit, ok := ret.Results[0].(*ast.CompositeLit); ok {
			funcBody = extractFuncBodyFromCompositeLit(compLit, pass.TypesInfo)
			if funcBody != nil {
				return false
			}
		}

		// Case 2 & 3: Function call
		if callExpr, ok := ret.Results[0].(*ast.CallExpr); ok {
			// Case 2: r.base.updateFunc() - method call
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if methodObj := pass.TypesInfo.Uses[selExpr.Sel]; methodObj != nil {
					if baseFunc := FindFuncDecl(pass, methodObj); baseFunc != nil {
						funcBody = GetFuncBody(pass, baseFunc)
						return false
					}
				}
			}

			// Case 3: someHelper() - standalone function call
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				funcObj := pass.TypesInfo.Uses[ident]
				if helperDecl := FindFuncDecl(pass, funcObj); helperDecl != nil {
					funcBody = GetFuncBody(pass, helperDecl)
					return false
				}
			}
		}

		return true
	})

	return funcBody
}

// Get schema map returned from func directly
func GetSchemaMapReturnedFromFunc(pass *analysis.Pass, funcDecl *ast.FuncDecl) *ast.CompositeLit {
	var schemaMap *ast.CompositeLit
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			return true
		}

		// Case 1: Direct return of composite literal
		if compLit, ok := ret.Results[0].(*ast.CompositeLit); ok {
			if IsSchemaMap(compLit, pass.TypesInfo) {
				schemaMap = compLit
				return false
			}
		}

		// Case 2: Return of a variable reference, only look into initial definition, ignoring later assignments
		if ident, ok := ret.Results[0].(*ast.Ident); ok {
			if compLit := TraceIdentToCompositeLit(pass.TypesInfo, ident, funcDecl); compLit != nil {
				if IsSchemaMap(compLit, pass.TypesInfo) {
					schemaMap = compLit
					return false
				}
			}
		}

		// Case 3: Return from function call
		if callExpr, ok := ret.Results[0].(*ast.CallExpr); ok {
			// Case 3a: r.base.arguments(schema) - method call with merge
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if methodObj := pass.TypesInfo.Uses[selExpr.Sel]; methodObj != nil {
					if baseFunc := FindFuncDecl(pass, methodObj); baseFunc != nil {
						baseSchemaMap := GetSchemaMapReturnedFromFunc(pass, baseFunc)
						if baseSchemaMap != nil && len(callExpr.Args) > 0 {
							if childMap := extractSchemaFromArg(pass, callExpr.Args[0], funcDecl); childMap != nil {
								schemaMap = mergeSchemaMaps(baseSchemaMap, childMap)
								return false
							}
						}
					}
				}
			}

			// Case 3b: getDeploymentScriptArguments() - standalone helper function
			if ident, ok := callExpr.Fun.(*ast.Ident); ok {
				funcObj := pass.TypesInfo.Uses[ident]
				if funcDecl := FindFuncDecl(pass, funcObj); funcDecl != nil {
					schemaMap = GetSchemaMapReturnedFromFunc(pass, funcDecl)
					return false
				}
			}
		}

		return true
	})

	return schemaMap
}

// TraceIdentToCompositeLit traces an identifier back to its first definition and returns the CompositeLit if found.
func TraceIdentToCompositeLit(typesInfo *types.Info, ident *ast.Ident, funcDecl *ast.FuncDecl) *ast.CompositeLit {
	obj := typesInfo.Uses[ident]
	if obj == nil {
		return nil
	}

	// Find the first definition of this variable
	var defNode ast.Node
	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		if assign, ok := node.(*ast.AssignStmt); ok {
			// Only check initial definitions (:=), not reassignments (=)
			if assign.Tok == token.DEFINE {
				for i, lhs := range assign.Lhs {
					if lhsIdent, ok := lhs.(*ast.Ident); ok {
						if typesInfo.Defs[lhsIdent] == obj && i < len(assign.Rhs) {
							defNode = assign.Rhs[i]
							return false
						}
					}
				}
			}
		}
		return defNode == nil
	})

	// Check if the definition is a composite literal
	if defNode != nil {
		if compLit, ok := defNode.(*ast.CompositeLit); ok {
			return compLit
		}
	}

	return nil
}

// IsTypedSDKResource checks if the type is from internal sdk pkg
func IsTypedSDKResource(t types.Type, name string) bool {
	switch t := t.(type) {
	case *types.Named:
		return IsNamedSDKResource(t, name)
	case *types.Pointer:
		return IsTypedSDKResource(t.Elem(), name)
	default:
		return false
	}
}

// IsNamedSDKResource checks if the named type is from internal sdk pkg
func IsNamedSDKResource(t *types.Named, name string) bool {
	obj := t.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == PackagePathSDK &&
		strings.HasPrefix(obj.Name(), name)
}

// extractSchemaFromArg extracts the schema map from a function argument
func extractSchemaFromArg(pass *analysis.Pass, arg ast.Expr, funcDecl *ast.FuncDecl) *ast.CompositeLit {
	switch v := arg.(type) {
	case *ast.CompositeLit:
		if IsSchemaMap(v, pass.TypesInfo) {
			return v
		}
	case *ast.Ident:
		if compLit := TraceIdentToCompositeLit(pass.TypesInfo, v, funcDecl); compLit != nil {
			if IsSchemaMap(compLit, pass.TypesInfo) {
				return compLit
			}
		}
	}
	return nil
}

// mergeSchemaMaps merges base schema map with child schema map (child overrides base)
func mergeSchemaMaps(baseMap, childMap *ast.CompositeLit) *ast.CompositeLit {
	if baseMap == nil {
		return childMap
	}
	if childMap == nil {
		return baseMap
	}

	// Build child field names set
	childFields := make(map[string]bool)
	for _, elt := range childMap.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.BasicLit); ok {
				childFields[strings.Trim(key.Value, `"`)] = true
			}
		}
	}

	// Merge: base fields (not overridden) + all child fields
	merged := &ast.CompositeLit{
		Type:   baseMap.Type,
		Lbrace: baseMap.Lbrace,
		Rbrace: childMap.Rbrace,
	}

	for _, elt := range baseMap.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.BasicLit); ok {
				if !childFields[strings.Trim(key.Value, `"`)] {
					merged.Elts = append(merged.Elts, elt)
				}
			}
		}
	}
	merged.Elts = append(merged.Elts, childMap.Elts...)

	return merged
}

// IsModelType checks if expression has the model type (handles: model, &model, *model, type aliases)
func IsModelType(expr ast.Expr, modelTypeName string, typesInfo *types.Info) bool {
	// Unwrap &model or *model at AST level
	if unary, ok := expr.(*ast.UnaryExpr); ok && (unary.Op == token.AND || unary.Op == token.MUL) {
		expr = unary.X
	}
	typ := typesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}
	return isModelTypeRecursive(typ, modelTypeName)
}

// isModelTypeRecursive checks type recursively
func isModelTypeRecursive(t types.Type, modelTypeName string) bool {
	switch t := t.(type) {
	case *types.Alias:
		return isModelTypeRecursive(types.Unalias(t), modelTypeName)
	case *types.Pointer:
		return isModelTypeRecursive(t.Elem(), modelTypeName)
	case *types.Named:
		if obj := t.Obj(); obj != nil {
			return obj.Name() == modelTypeName
		}
	}
	return false
}
