package helper

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/bflad/tfproviderlint/helper/astutils"
)

const (
	TypeNameResourceData = "ResourceData"
)

func IsResourceData(typesInfo *types.Info, sel *ast.SelectorExpr) bool {
	if typesInfo == nil {
		return false
	}

	// Get the type of the selector's base expression
	typ := typesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}

	return IsTypeResourceData(typ)
}

func IsTypeResourceData(t types.Type) bool {
	switch t := t.(type) {
	case *types.Alias:
		return IsTypeResourceData(types.Unalias(t))
	case *types.Named:
		return astutils.IsModulePackageNamedType(t, ModuleTerraformPluginSDK, PackageModulePathSchema, TypeNameResourceData) ||
			IsNamedResourceData(t)
	case *types.Pointer:
		return IsTypeResourceData(t.Elem())
	default:
		return false
	}
}

func IsNamedResourceData(t *types.Named) bool {
	obj := t.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return (obj.Pkg().Path() == PackagePathSDK || obj.Pkg().Path() == PackagePathPluginSDK) &&
		strings.HasPrefix(obj.Name(), TypeNameResourceData)
}
