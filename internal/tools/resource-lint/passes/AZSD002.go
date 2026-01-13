package passes

import (
	"go/ast"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	localschema "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
)

const AZSD002Doc = `check AtLeastOneOf validation for TypeList fields with all optional nested fields

The AZSD002 analyzer checks that when a pluginsdk.TypeList block has no required nested 
fields, AtLeastOneOf must be set on the optional fields to ensure at least one is specified.

Example violation:
  "setting": {
      Type:     pluginsdk.TypeList,
      Optional: true,
      MaxItems: 1,
      Elem: &pluginsdk.Resource{
          Schema: map[string]*pluginsdk.Schema{
              "linux": {
                  Type:     pluginsdk.TypeList,
                  Optional: true,
                  // Missing AtLeastOneOf!
              },
              "windows": {
                  Type:     pluginsdk.TypeList,
                  Optional: true,
                  // Missing AtLeastOneOf!
              },
          },
      },
  }

Valid usage:
  "setting": {
      Type:     pluginsdk.TypeList,
      Optional: true,
      MaxItems: 1,
      Elem: &pluginsdk.Resource{
          Schema: map[string]*pluginsdk.Schema{
              "linux": {
                  Type:         pluginsdk.TypeList,
                  Optional:     true,
                  AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
              },
              "windows": {
                  Type:         pluginsdk.TypeList,
                  Optional:     true,
                  AtLeastOneOf: []string{"setting.0.linux", "setting.0.windows"},
              },
          },
      },
  }`

const azsd002Name = "AZSD002"

var AZSD002Analyzer = &analysis.Analyzer{
	Name:     azsd002Name,
	Doc:      AZSD002Doc,
	Run:      runAZSD002,
	Requires: []*analysis.Analyzer{localschema.LocalAnalyzer},
}

func runAZSD002(pass *analysis.Pass) (interface{}, error) {
	schemaInfoCache, ok := pass.ResultOf[localschema.LocalAnalyzer].(map[*ast.CompositeLit]*localschema.LocalSchemaInfoWithName)
	if !ok {
		return nil, nil
	}

	for schemaLit, cached := range schemaInfoCache {
		schemaInfo := cached.Info

		// Skip Computed fields
		if cached.Info.Schema.Computed {
			continue
		}

		// Only check TypeList fields
		if !schemaInfo.IsType(schema.SchemaValueTypeList) {
			continue
		}

		// Get Elem field
		elemKV := schemaInfo.Fields[schema.SchemaFieldElem]
		if elemKV == nil {
			continue
		}

		// Check if Elem is &schema.Resource{...}
		resourceSchema := helper.GetResourceSchemaFromElem(elemKV)
		if resourceSchema == nil {
			continue
		}

		// Find the Schema field in the Resource
		nestedSchemaMap := helper.GetNestedSchemaMap(resourceSchema)
		if nestedSchemaMap == nil {
			continue
		}

		// Collect nested fields
		optionalFieldsCount := 0
		hasRequiredField := false
		hasAtLeastOneOf := false
		for _, elt := range nestedSchemaMap.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			nestedSchemaLit, ok := kv.Value.(*ast.CompositeLit)
			if !ok {
				continue
			}

			nestedCached, exists := schemaInfoCache[nestedSchemaLit]
			if !exists {
				continue
			}

			nestedInfo := nestedCached.Info
			if nestedInfo.Schema.Required {
				hasRequiredField = true
				break
			}

			if nestedInfo.Schema.Optional {
				// Check if at least one optional field has AtLeastOneOf
				atLeastOneOfKV := nestedInfo.Fields[schema.SchemaFieldAtLeastOneOf]
				if atLeastOneOfKV != nil {
					hasAtLeastOneOf = true
					break
				}
				optionalFieldsCount++
			}
		}

		// Only report if there are no required fields, multiple optional fields,
		// and none of them have AtLeastOneOf set
		if !hasRequiredField && !hasAtLeastOneOf && optionalFieldsCount >= 2 {
			pos := pass.Fset.Position(schemaLit.Pos())
			if loader.ShouldReport(pos.Filename, pos.Line) {
				if propertyName := cached.PropertyName; propertyName != "" {
					pass.Reportf(schemaLit.Pos(),
						"%s: TypeList field `%s` has %s, %s must be set on the optional fields to ensure at least one is specified.\n",
						azsd002Name, propertyName, helper.IssueLine("all optional nested fields"), helper.FixedCode("`AtLeastOneOf`"))
				} else {
					pass.Reportf(schemaLit.Pos(),
						"%s: TypeList field has %s, %s must be set on the optional fields to ensure at least one is specified.\n",
						azsd002Name, helper.IssueLine("all optional nested fields"), helper.FixedCode("`AtLeastOneOf`"))
				}
			}
		}
	}

	return nil, nil
}
