package helper

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/bflad/tfproviderlint/helper/astutils"
	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"golang.org/x/tools/go/analysis"
)

const (
	TypeNameSchema = "Schema"

	// Package paths for schema types
	ModuleTerraformPluginSDK = "github.com/hashicorp/terraform-plugin-sdk/v2"
	PackageModulePathSchema  = "helper/schema"

	// azurerm provider specific package
	PackagePathPluginSDK = "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// SchemaFieldInfo represents a field in a Terraform schema with its schema information
type SchemaFieldInfo struct {
	Name       string
	SchemaInfo *schema.SchemaInfo
	Position   int
}

// IsSchemaMap checks if a composite literal is a map[string]*schema.Schema or map[string]*pluginsdk.Schema
func IsSchemaMap(cl *ast.CompositeLit, info *types.Info) bool {
    // Check if it's a map literal
    mapType, ok := cl.Type.(*ast.MapType)
    if !ok {
        return false
    }

    // Check if key is string
    switch k := mapType.Key.(type) {
    case *ast.Ident:
        if k.Name != "string" {
            return false
        }
    default:
        return false
    }

    return isTypeSchema(info.TypeOf(mapType.Value))
}

// IsSchemaSchema checks if a composite literal is of type schema.Schema or pluginsdk.Schema
func IsSchemaSchema(typesInfo *types.Info, cl *ast.CompositeLit) bool {
	if cl.Type == nil {
		return false
	}

	t := typesInfo.TypeOf(cl.Type)
	if t == nil {
		return false
	}

	return isTypeSchema(t)
}

// isTypeSchema returns if the type is Schema from helper/schema or pluginsdk package
func isTypeSchema(t types.Type) bool {
	switch t := t.(type) {
	case *types.Alias:
		return isTypeSchema(types.Unalias(t))
	case *types.Named:
		// Check if it's from helper/schema package
		if astutils.IsModulePackageNamedType(t, ModuleTerraformPluginSDK, PackageModulePathSchema, TypeNameSchema) {
			return true
		}
		// Check if it's from pluginsdk package (azurerm provider specific)
		if t.Obj().Name() == TypeNameSchema && t.Obj().Pkg() != nil &&
			t.Obj().Pkg().Path() == PackagePathPluginSDK {
			return true
		}
		return false
	case *types.Pointer:
		return isTypeSchema(t.Elem())
	default:
		return false
	}
}

// IsNestedSchemaMap checks if a schema map CompositeLit is nested within an Elem field
func IsNestedSchemaMap(file *ast.File, schemaLit *ast.CompositeLit) bool {
	var isNested bool

	ast.Inspect(file, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}

		// Check if this is an Elem key
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Elem" {
			return true
		}

		// Check if our schemaLit is within this Elem value's range
		if schemaLit.Pos() >= kv.Value.Pos() && schemaLit.End() <= kv.Value.End() {
			isNested = true
			return false // Found it, stop searching immediately
		}

		return true
	})

	return isNested
}

// FindFuncDecl finds the function declaration for a given function object
func FindFuncDecl(pass *analysis.Pass, funcObj types.Object) *ast.FuncDecl {
	obj, ok := funcObj.(*types.Func)
	if !ok {
		return nil
	}

	pos := obj.Pos()

	for _, file := range pass.Files {
		// Check if the position is within this file's range
		if pos < file.Pos() || pos > file.End() {
			continue
		}

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Match by function name position
			if funcDecl.Name.Pos() == pos {
				return funcDecl
			}
		}
	}

	return nil
}

// GetResourceSchemaFromElem extracts the &schema.Resource{...} composite literal from an Elem field
// Returns nil if Elem is not a pointer to a Resource composite literal
func GetResourceSchemaFromElem(elemKV *ast.KeyValueExpr) *ast.CompositeLit {
	if unary, ok := elemKV.Value.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		if compLit, ok := unary.X.(*ast.CompositeLit); ok {
			return compLit
		}
	}

	return nil
}

// GetNestedSchemaMap extracts the Schema field value from a &schema.Resource{...} composite literal
// Returns nil if the Schema field is not found or is not a composite literal
func GetNestedSchemaMap(resourceSchema *ast.CompositeLit) *ast.CompositeLit {
    fields := astutils.CompositeLitFields(resourceSchema)
    if kvExpr := fields[schema.ResourceFieldSchema]; kvExpr != nil {
        if compLit, ok := kvExpr.Value.(*ast.CompositeLit); ok {
            return compLit
        }
    }
    return nil
}
