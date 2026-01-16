package passes

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strings"

	"github.com/bflad/tfproviderlint/helper/astutils"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"

	"golang.org/x/tools/go/analysis"
)

const AZNR002Doc = `check that top-level updatable properties are handled in Update function

The AZNR002 analyzer checks that all updatable properties (not marked as ForceNew)
are properly handled in the Update function for typed resources.

If git filter enabled, this rule only applies if schema is changed.

For typed resources, this means checking for metadata.ResourceData.HasChange("property_name").

Note: This analyzer supports Arguments() functions that:
 - Directly return map[string]*pluginsdk.Schema{}
 - Return a variable (traces to initial := definition, ignoring subsequent modifications)

Example violation:
  // In Arguments()
  "display_name": {
      Type:     pluginsdk.TypeString,
      Required: true,
      // No ForceNew - this is updatable
  }

  // In Update() - missing HasChange check
  func (r Resource) Update() sdk.ResourceFunc {
      return sdk.ResourceFunc{
          Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
              // Missing: if metadata.ResourceData.HasChange("display_name") { ... }
              return nil
          },
      }
  }

Valid usage:
  func (r Resource) Update() sdk.ResourceFunc {
      return sdk.ResourceFunc{
          Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
              if metadata.ResourceData.HasChange("display_name") {
                  props.DisplayName = pointer.To(config.DisplayName)
              }
              return nil
          },
      }
  }`

const aznr002Name = "AZNR002"

var aznr002SkipPackages = []string{"_test", "/migration", "/client", "/validate", "/test-data", "/parse", "/models"}

var AZNR002Analyzer = &analysis.Analyzer{
	Name:     aznr002Name,
	Doc:      AZNR002Doc,
	Run:      runAZNR002,
	Requires: []*analysis.Analyzer{schema.TypedResourceInfoAnalyzer},
}

func runAZNR002(pass *analysis.Pass) (interface{}, error) {
	pkgPath := pass.Pkg.Path()
	for _, skip := range aznr002SkipPackages {
		if strings.Contains(pkgPath, skip) {
			return nil, nil
		}
	}

	allTypedResources, ok := pass.ResultOf[schema.TypedResourceInfoAnalyzer].([]*helper.TypedResourceInfo)
	if !ok {
		return nil, nil
	}
	for _, resource := range allTypedResources {
		// Filter: must have Update method
		if resource.UpdateFunc == nil {
			continue
		}

		// Filter: must have extracted ArgumentsProperties
		if len(resource.ArgumentsProperties) == 0 {
			pos := pass.Fset.Position(resource.ArgumentsFunc.Pos())
			log.Printf("%s:%d: %s: Skipping resource %q - failed to extract schema properties",
				pos.Filename, pos.Line, aznr002Name, resource.ResourceTypeName)
			continue
		}

		updatableProps := extractUpdatableProperties(resource)
		if len(updatableProps) == 0 {
			continue
		}

		handledProps := findHandledPropertiesInUpdate(pass, resource)
		reportMissingProperties(pass, resource, updatableProps, handledProps)
	}

	return nil, nil
}

// extractUpdatableProperties filters updatable properties from ArgumentsProperties
func extractUpdatableProperties(resource *helper.TypedResourceInfo) map[string]string {
	updatableProps := make(map[string]string)

	tfSchemaToModel := make(map[string]string)
	for modelField, tfSchema := range resource.ModelFieldToTFSchema {
		tfSchemaToModel[tfSchema] = modelField
	}

	for _, field := range resource.ArgumentsProperties {
		if field.SchemaInfo != nil &&
			!field.SchemaInfo.Schema.Computed &&
			!field.SchemaInfo.Schema.ForceNew {
			updatableProps[field.Name] = tfSchemaToModel[field.Name]
		}
	}

	return updatableProps
}

// findHandledPropertiesInUpdate finds all properties handled in Update function
func findHandledPropertiesInUpdate(pass *analysis.Pass, resource *helper.TypedResourceInfo) map[string]bool {
	handledProps := make(map[string]bool)

	updateFuncBody := resource.UpdateFuncBody
	if updateFuncBody == nil {
		return handledProps
	}

	// Get the model struct type name
	modelTypeName := resource.ModelName

	// Pattern 1: Trace into helper functions that receive model or metadata (recursive)
	traceHelperCalls(pass, updateFuncBody, modelTypeName, resource, handledProps, make(map[*ast.FuncDecl]bool))

	// Pattern 2: Direct model field access in Update body (config.Field, state.Field)
	traceModelFieldAccess(updateFuncBody, modelTypeName, resource, handledProps)

	// Pattern 3: ResourceData method calls (HasChange/HasChanges/Get)
	traceResourceDataCalls(updateFuncBody, resource, handledProps)

	return handledProps
}

// traceHelperCalls recursively traces into helper functions that receive model or metadata as argument
// Uses type-based detection to find relevant variables
// visited tracks already-visited functions to prevent infinite recursion
func traceHelperCalls(pass *analysis.Pass, body *ast.BlockStmt, modelTypeName string, resource *helper.TypedResourceInfo, handledProps map[string]bool, visited map[*ast.FuncDecl]bool) {
	if body == nil {
		return
	}
	typesInfo := resource.TypesInfo

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Skip metadata.Decode()
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "Decode" && helper.IsTypedSDKResource(typesInfo.TypeOf(sel.X), "ResourceMetaData") {
				return true
			}
		}

		// Check if model or metadata is passed as an argument
		shouldTrace := false
		for _, arg := range call.Args {
			if helper.IsModelType(arg, modelTypeName, typesInfo) || helper.IsTypedSDKResource(typesInfo.TypeOf(arg), "ResourceMetaData") {
				shouldTrace = true
				break
			}
		}
		if !shouldTrace {
			return true
		}

		// Resolve and trace into the helper function
		funcDecl := resolveFuncDecl(pass, call)
		if funcDecl == nil || funcDecl.Body == nil || visited[funcDecl] {
			return true
		}
		visited[funcDecl] = true

		// Trace all patterns in helper body
		traceModelFieldAccess(funcDecl.Body, modelTypeName, resource, handledProps)
		traceResourceDataCalls(funcDecl.Body, resource, handledProps)
		traceHelperCalls(pass, funcDecl.Body, modelTypeName, resource, handledProps, visited)

		return true
	})
}

// resolveFuncDecl resolves a CallExpr to its FuncDecl (same package only)
func resolveFuncDecl(pass *analysis.Pass, call *ast.CallExpr) *ast.FuncDecl {
	var funcObj types.Object
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		funcObj = pass.TypesInfo.Uses[fun]
	case *ast.SelectorExpr:
		funcObj = pass.TypesInfo.Uses[fun.Sel]
	}
	if funcObj == nil {
		return nil
	}
	return helper.FindFuncDecl(pass, funcObj)
}

// traceModelFieldAccess finds model.Field accesses in a function body using type detection
func traceModelFieldAccess(body *ast.BlockStmt, modelTypeName string, resource *helper.TypedResourceInfo, handledProps map[string]bool) {
	if body == nil || resource.TypesInfo == nil {
		return
	}

	ast.Inspect(body, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Check if field name maps to a tfschema property
		tfschemaName, ok := resource.ModelFieldToTFSchema[sel.Sel.Name]
		if !ok {
			return true
		}

		// Verify the base is a model variable by type
		if helper.IsModelType(sel.X, modelTypeName, resource.TypesInfo) {
			handledProps[tfschemaName] = true
		}
		return true
	})
}

// traceResourceDataCalls finds ResourceData method calls (HasChange/HasChanges/Get) in a function body
func traceResourceDataCalls(body *ast.BlockStmt, resource *helper.TypedResourceInfo, handledProps map[string]bool) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if helper.IsResourceData(resource.TypesInfo, sel) {
			for _, arg := range call.Args {
				if propName := astutils.ExprStringValue(arg); propName != nil {
					handledProps[*propName] = true
				}
			}
		}
		return true
	})
}

// reportMissingProperties reports properties that are updatable but not handled
func reportMissingProperties(pass *analysis.Pass, resource *helper.TypedResourceInfo, updatableProps map[string]string, handledProps map[string]bool) {
	var missingProps []string
	for tfSchemaName := range updatableProps {
		if !handledProps[tfSchemaName] {
			missingProps = append(missingProps, tfSchemaName)
		}
	}

	if len(missingProps) == 0 || len(handledProps) == 0 {
		if len(handledProps) == 0 && len(updatableProps) != 0 {
			pos := pass.Fset.Position(resource.UpdateFunc.Pos())
			log.Printf("%s:%d: %s: Skipping resource %q - update likely delegated to helper function",
				pos.Filename, pos.Line, aznr002Name, resource.ResourceTypeName)
		}
		return
	}

	// Sort for consistent output
	sort.Strings(missingProps)

	// Report each missing property at its definition location
	for _, tfSchemaName := range missingProps {
		var fieldInfo *helper.SchemaFieldInfo
		for i := range resource.ArgumentsProperties {
			if resource.ArgumentsProperties[i].Name == tfSchemaName {
				fieldInfo = &resource.ArgumentsProperties[i]
				break
			}
		}
		if fieldInfo == nil {
			continue
		}

		if fieldInfo.Pos != token.NoPos {
			position := pass.Fset.Position(fieldInfo.Pos)
			// Check if position is valid (Pos is in current pass's FileSet)
			if position.IsValid() {
				if !loader.ShouldReport(position.Filename, position.Line) {
					continue
				}
				pass.Reportf(fieldInfo.Pos,
					"%s: updatable property `%s` is not handled in Update function. If non-updatable, mark as %s in Arguments() schema\n",
					aznr002Name,
					helper.IssueLine(tfSchemaName),
					helper.FixedCode("ForceNew: true"))
				continue
			}
		}

		// Fallback to Update function position (for cross-package schemas)
		if fieldInfo.Position.IsValid() {
			if !loader.ShouldReport(fieldInfo.Position.Filename, fieldInfo.Position.Line) {
				continue
			}
		}
		pass.Reportf(resource.UpdateFunc.Pos(),
			"%s: updatable property `%s` is not handled in Update function. If non-updatable, mark as %s in Arguments() schema\n",
			aznr002Name,
			helper.IssueLine(tfSchemaName),
			helper.FixedCode("ForceNew: true"))
	}
}
