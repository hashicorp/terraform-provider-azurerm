package passes

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const PointerPkgPath = "github.com/hashicorp/go-azure-helpers/lang/pointer"
const azbp003Name = "AZBP003"

const AZBP003Doc = `check enum conversion for go-azure-sdk

The AZBP003 analyzer checks that when enum types are converted using pointer.To() with explicit type conversion instead of pointer.ToEnum[]

Example violation:
  return &managedclusters.ManagedClusterBootstrapProfile{
    ArtifactSource: pointer.To(managedclusters.ArtifactSource(config["artifact_source"].(string))),
  }

Valid usage:
  return &managedclusters.ManagedClusterBootstrapProfile{
    ArtifactSource: pointer.ToEnum[managedclusters.ArtifactSource](config["artifact_source"].(string)),
  }`

var AZBP003Analyzer = &analysis.Analyzer{
	Name:     azbp003Name,
	Doc:      AZBP003Doc,
	Run:      runAZBP003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var azbp003SkipPackages = []string{"_test", "/migration", "/client", "/validate", "/test-data", "/parse", "/models"}

func runAZBP003(pass *analysis.Pass) (interface{}, error) {
	// Skip specified packages
	pkgPath := pass.Pkg.Path()
	for _, skip := range azbp003SkipPackages {
		if strings.Contains(pkgPath, skip) {
			return nil, nil
		}
	}

	relevantFiles := make(map[string]bool)
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		importsPointer := false
		for _, imp := range f.Imports {
			if imp.Path.Value == `"`+PointerPkgPath+`"` {
				importsPointer = true
				break
			}
		}

		if importsPointer {
			relevantFiles[filename] = true
		}
	}

	if len(relevantFiles) == 0 {
		return nil, nil
	}

	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	// Pre-filter: only look at CallExpr nodes
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return
		}

		// pattern: pointer.To(sdk.EnumConstant(xxx(string)))
		argCall, ok := call.Args[0].(*ast.CallExpr)
		if !ok {
			return
		}

		pos := pass.Fset.Position(call.Pos())
		filename := pos.Filename
		if !relevantFiles[filename] {
			return
		}

		// selector expressions: e.g. pointer.To
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "To" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "pointer" {
			return
		}

		// Validate pkg path
		obj := pass.TypesInfo.Uses[ident]
		if obj == nil {
			return
		}

		pkg, ok := obj.(*types.PkgName)
		if !ok || pkg.Imported().Path() != PointerPkgPath {
			return
		}

		targetType := pass.TypesInfo.TypeOf(argCall.Fun)
		if targetType == nil {
			return
		}

		named, ok := targetType.(*types.Named)
		if !ok || !isEnumTypeInSDK(pass, named) {
			return
		}

		if loader.ShouldReport(pos.Filename, pos.Line) {
			pass.Reportf(call.Pos(), "%s: use `%s` to convert Enum type instead of explicitly type conversion.\n",
				azbp003Name, helper.FixedCode("pointer.ToEnum"))
		}
	})

	return nil, nil
}

// Find if it's an enum type by checking for PossibleValuesFor{TypeName} function
func isEnumTypeInSDK(pass *analysis.Pass, named *types.Named) bool {
	// 1. Check if underlying type is string OR integer
	basic, ok := named.Underlying().(*types.Basic)
	if !ok {
		return false
	}

	// Accept string or integer types (int, int64, int32, etc.)
	info := basic.Info()
	if info&types.IsString == 0 && info&types.IsInteger == 0 {
		return false
	}

	// 2. Check package path is go Azure SDK
	pkg := named.Obj().Pkg()
	if pkg == nil || !strings.Contains(pkg.Path(), "github.com/hashicorp/go-azure-sdk") {
		return false
	}

	// 3. Check for PossibleValuesFor{TypeName} function - the standard enum pattern
	typeName := named.Obj().Name()
	functionName := "PossibleValuesFor" + typeName

	obj := pkg.Scope().Lookup(functionName)
	if obj == nil {
		// Fallback: check if defined in constants.go
		pos := named.Obj().Pos()
		position := pass.Fset.Position(pos)
		return strings.HasSuffix(position.Filename, "constants.go")
	}

	// Verify it's a function returning []string
	fn, ok := obj.(*types.Func)
	if !ok {
		return false
	}

	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}
	if sig.Params().Len() != 0 || sig.Results().Len() != 1 {
		return false
	}

	// Check return type is []string
	slice, ok := sig.Results().At(0).Type().(*types.Slice)
	if !ok {
		return false
	}

	elem, ok := slice.Elem().(*types.Basic)
	if !ok {
		return false
	}

	return ok && elem.Kind() == basic.Kind()
}
