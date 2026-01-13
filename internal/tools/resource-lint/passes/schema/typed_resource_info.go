package schema

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const typedResourceInfoDoc = `Finds all typed resources and extracts their schema information.

Key Features:
 1. Identifies SDK resource interface declarations (var _ sdk.Resource = &MyResource{})
 2. Extracts Arguments() function and its returned schema map
 3. Handles both direct returns and variable returns (traces to initial := definition only)
 4. Parses schema properties using ExtractSchemaInfoFromMap with commonschema support
 5. Deduplicates resources that implement multiple SDK interfaces

Example:

	// Resource declaration
	var _ sdk.Resource = &MyResource{}
	var _ sdk.ResourceWithUpdate = &MyResource{}  // Won't be processed twice

	type MyResource struct{}

	// Case 1: Direct return
	func (r MyResource) Arguments() map[string]*pluginsdk.Schema {
	    return map[string]*pluginsdk.Schema{
	        "name": {Type: TypeString, Required: true},
	        "resource_group_name": commonschema.ResourceGroupName(),
	    }
	}

	// Case 2: Variable return (captures initial definition only)
	func (r AnotherResource) Arguments() map[string]*pluginsdk.Schema {
	    output := map[string]*pluginsdk.Schema{  // <- Captures this
	        "name": {Type: TypeString, Required: true},
	    }
	    output["tags"] = Tags()  // <- Ignores modifications
	    return output
	}

Processing:
 - Only processes files ending with _resource.go
 - Uses type system to check SDK interface implementations
 - Traces variable returns to initial := definition, ignoring subsequent = assignments
`

var TypedResourceInfoAnalyzer = &analysis.Analyzer{
	Name: "typedresourceinfo",
	Doc:  typedResourceInfoDoc,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		CompleteSchemaAnalyzer,
	},
	Run:        runTypedResourceInfo,
	ResultType: reflect.TypeOf([]*helper.TypedResourceInfo{}),
}

func runTypedResourceInfo(pass *analysis.Pass) (interface{}, error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}
	completeSchemaInfo, ok := pass.ResultOf[CompleteSchemaAnalyzer].(*CompleteSchemaInfo)
	if !ok {
		return nil, nil
	}

	var result []*helper.TypedResourceInfo
	seen := make(map[string]bool)

	nodeFilter := []ast.Node{(*ast.GenDecl)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return
		}

		fileName := pass.Fset.Position(genDecl.Pos()).Filename
		if !strings.HasSuffix(fileName, "_resource.go") {
			return
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if !isSDKResourceInterface(pass, valueSpec.Type) {
				continue
			}

			if len(valueSpec.Values) == 0 {
				continue
			}

			var resourceTypeName string
			if compLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					resourceTypeName = ident.Name
				}
			}

			if resourceTypeName == "" {
				continue
			}

			// Skip if already processed
			if seen[resourceTypeName] {
				continue
			}
			seen[resourceTypeName] = true

			for _, file := range pass.Files {
				if pass.Fset.Position(file.Pos()).Filename != fileName {
					continue
				}

				resourceInfo := helper.NewTypedResourceInfo(resourceTypeName, file, pass.TypesInfo)
				if resourceInfo.ArgumentsFunc == nil {
					continue
				}

				schemaMap := helper.GetSchemaMapReturnedFromFunc(pass, resourceInfo.ArgumentsFunc)
				if schemaMap == nil {
					result = append(result, resourceInfo)
					continue
				}

				resourceInfo.ArgumentsProperties = completeSchemaInfo.SchemaFields[schemaMap.Pos()]

				result = append(result, resourceInfo)
			}
		}
	})

	return result, nil
}

// isSDKResourceInterface checks if the type is an sdk.Resource* interface
func isSDKResourceInterface(pass *analysis.Pass, expr ast.Expr) bool {
	selExpr, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	typ := pass.TypesInfo.TypeOf(selExpr)
	if typ != nil && helper.IsTypedSDKResource(typ, "Resource") {
		return true
	}

	// Fallback: AST-based check for testdata/mock packages
	pkgIdent, ok := selExpr.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "sdk" {
		return false
	}

	return strings.HasPrefix(selExpr.Sel.Name, "Resource")
}
