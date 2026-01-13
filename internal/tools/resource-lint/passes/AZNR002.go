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

If git filter enabled, this rule only applies on newly created typed resource.

This analyzer will be skipped if a helper function is utilized to handle the update.

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

	allResources, ok := pass.ResultOf[schema.TypedResourceInfoAnalyzer].([]*helper.TypedResourceInfo)
	if !ok {
		return nil, nil
	}
	for _, resource := range allResources {
		// Filter: must have Update method
		if resource.UpdateFunc == nil {
			continue
		}

		// Filter: git filter - only check new files
		fileName := pass.Fset.Position(resource.UpdateFunc.Pos()).Filename
		if !loader.IsNewFile(fileName) {
			continue
		}

		// Filter: must have extracted ArgumentsProperties
		if len(resource.ArgumentsProperties) == 0 {
			pos := pass.Fset.Position(resource.UpdateFunc.Pos())
			log.Printf("%s:%d: %s: Skipping resource %q - failed to extract schema properties",
				pos.Filename, pos.Line, aznr002Name, resource.ResourceTypeName)
			continue
		}

		updatableProps := extractUpdatableProperties(resource)
		if len(updatableProps) == 0 {
			continue
		}

		handledProps := findHandledPropertiesInUpdate(resource)
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
func findHandledPropertiesInUpdate(resource *helper.TypedResourceInfo) map[string]bool {
	handledProps := make(map[string]bool)

	updateFuncBody := resource.UpdateFuncBody
	if updateFuncBody == nil {
		return handledProps
	}

	// Get the model struct type name
	modelTypeName := resource.ModelName

	// Pattern 1: Check if model/config is passed to helper functions
	// If detected, skip this resource as properties are likely handled in helper
	if detectModelPassedToHelper(updateFuncBody, modelTypeName, resource.TypesInfo) {
		return handledProps
	}

	// Single pass: inspect all nodes and check both patterns
	ast.Inspect(updateFuncBody, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name

				// Pattern 2 & 3: Check ResourceData method calls (HasChange/HasChanges/Get)
				if methodName == "HasChange" || methodName == "HasChanges" || methodName == "Get" {
					if helper.IsResourceData(resource.TypesInfo, sel) {
						if methodName == "Get" && len(node.Args) > 0 {
							// Pattern 3: Get("property_name")
							if propName := astutils.ExprStringValue(node.Args[0]); propName != nil {
								handledProps[*propName] = true
							}
						} else if methodName == "HasChange" || methodName == "HasChanges" {
							// Pattern 2: HasChange("prop") or HasChanges("prop1", "prop2")
							for _, arg := range node.Args {
								if propName := astutils.ExprStringValue(arg); propName != nil {
									handledProps[*propName] = true
								}
							}
						}
					}
				}
			}

		case *ast.SelectorExpr:
			// Pattern 4: state.FieldName or config.FieldName
			// Check if the field name matches any of our model fields
			fieldName := node.Sel.Name
			if tfschemaName, ok := resource.ModelFieldToTFSchema[fieldName]; ok {
				// This is a field from our model struct being accessed
				// Now verify the base is likely a model variable by checking with TypesInfo
				if resource.TypesInfo != nil {
					if typ := resource.TypesInfo.TypeOf(node.X); typ != nil {
						// Remove pointer if present
						if ptr, ok := typ.(*types.Pointer); ok {
							typ = ptr.Elem()
						}
						// Check if it's a named type matching our model
						if named, ok := typ.(*types.Named); ok {
							if obj := named.Obj(); obj != nil && obj.Name() == modelTypeName {
								handledProps[tfschemaName] = true
							}
						}
					}
				}
			}
		}

		return true
	})

	return handledProps
}

// detectModelPassedToHelper checks if model/config variable is passed to helper functions
// Returns true if expand/map/flatten functions are called with model variables at TOP LEVEL (not inside if/for blocks)
// e.g. "automanage_configuration_resource.go": expandConfigurationProfile(model) - should skip
// counter-example "spring_cloud_gateway_resource.go": if HasChange { expandGatewayResponseCacheProperties(model) } - should NOT skip
func detectModelPassedToHelper(body *ast.BlockStmt, modelTypeName string, typesInfo *types.Info) bool {
	// Only check top-level statements in the function body
	for _, stmt := range body.List {
		// Skip if statement, for statement, switch statement - these are conditional
		switch stmt.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.RangeStmt:
			continue
		}

		found := false
		// Check assignments and expression statements at top level
		ast.Inspect(stmt, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Skip known SDK methods that don't delegate update logic
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				methodName := sel.Sel.Name
				// Skip metadata.Decode() - this is serialization, not business logic
				if methodName == "Decode" && helper.IsTypedSDKResource(typesInfo.TypeOf(sel.X), "ResourceMetaData") {
					return true
				}
			}

			// Check if any argument is the model variable (by type)
			for _, arg := range call.Args {
				var argIdent *ast.Ident
				// Handle: model, &model, *model
				switch a := arg.(type) {
				case *ast.Ident:
					argIdent = a
				case *ast.UnaryExpr:
					if a.Op == token.AND || a.Op == token.MUL {
						if ident, ok := a.X.(*ast.Ident); ok {
							argIdent = ident
						}
					}
				}

				if argIdent != nil {
					// Use TypesInfo to check if this variable is of model type
					if typ := typesInfo.TypeOf(argIdent); typ != nil {
						// Remove pointer if present to get underlying type
						if ptr, ok := typ.(*types.Pointer); ok {
							typ = ptr.Elem()
						}
						// Check if it's a named type matching our model
						if named, ok := typ.(*types.Named); ok {
							if obj := named.Obj(); obj != nil && obj.Name() == modelTypeName {
								found = true
								return false
							}
						}
					}
				}
			}

			return true
		})

		if found {
			return true
		}
	}

	return false
}

// reportMissingProperties reports properties that are updatable but not handled
func reportMissingProperties(pass *analysis.Pass, resource *helper.TypedResourceInfo, updatableProps map[string]string, handledProps map[string]bool) {
	var missingProps []string

	for propName := range updatableProps {
		if !handledProps[propName] {
			missingProps = append(missingProps, propName)
		}
	}

	if len(missingProps) == 0 || len(handledProps) == 0 {
		if len(handledProps) == 0 {
			pos := pass.Fset.Position(resource.UpdateFunc.Pos())
			log.Printf("%s:%d: %s: Skipping resource %q - update likely delegated to helper function",
				pos.Filename, pos.Line, aznr002Name, resource.ResourceTypeName)
		}
		return
	}

	// Sort for consistent output
	sort.Strings(missingProps)

	// Report at the Update function
	if resource.UpdateFunc != nil {
		pass.Reportf(resource.UpdateFunc.Pos(),
			"%s: resource has updatable properties not handled in Update function: `%s`. If they are non-updatable, mark them as %s in Arguments() schema\n",
			aznr002Name,
			helper.IssueLine(strings.Join(missingProps, ", ")),
			helper.FixedCode("ForceNew: true"))
	}
}
