package passes

import (
	"go/ast"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	localschema "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
)

const AZSD001Doc = `check MaxItems:1 blocks with single property should be flattened

The AZSD001 analyzer checks that blocks with MaxItems: 1 containing only a single 
nested property should be flattened unless there's a comment explaining why.

Example violation:
  "config": {
      Type:     schema.TypeList,
      MaxItems: 1,
      Elem: &schema.Resource{
          Schema: map[string]*schema.Schema{
              "value": {...},  // Only one property - should be flattened
          },
      },
  }

Valid usage (flattened):
  "config_value": {...}

Valid usage (with explanation):
  "config": {
      Type:     schema.TypeList,
      MaxItems: 1,
      // Additional properties will be added per service team confirmation
      Elem: &schema.Resource{
          Schema: map[string]*schema.Schema{
              "value": {...},
          },
      },
  }`

const azsd001Name = "AZSD001"

var AZSD001Analyzer = &analysis.Analyzer{
	Name:     azsd001Name,
	Doc:      AZSD001Doc,
	Run:      runAZSD001,
	Requires: []*analysis.Analyzer{localschema.LocalAnalyzer},
}

func runAZSD001(pass *analysis.Pass) (interface{}, error) {
	schemaInfoCache, ok := pass.ResultOf[localschema.LocalAnalyzer].(map[*ast.CompositeLit]*localschema.LocalSchemaInfoWithName)
	if !ok {
		return nil, nil
	}

	// Build file comments map for all files
	fileCommentsMap := make(map[string][]*ast.CommentGroup)
	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename
		fileCommentsMap[filename] = f.Comments
	}

	// Iterate over cached schema infos
	for schemaLit, cached := range schemaInfoCache {
		schemaInfo := cached.Info

		// Check if MaxItems is 1
		if schemaInfo.Schema.MaxItems != 1 {
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

		// Count properties in the nested schema
		propertyCount := 0
		for _, elt := range nestedSchemaMap.Elts {
			if _, ok := elt.(*ast.KeyValueExpr); ok {
				propertyCount++
			}
		}

		// If only one property, check for any explanatory comment
		if propertyCount == 1 {
			filename := pass.Fset.Position(schemaLit.Pos()).Filename
			elemLine := pass.Fset.Position(elemKV.Value.Pos()).Line

			hasComment := false
			comments := fileCommentsMap[filename]
			for _, cg := range comments {
				for _, c := range cg.List {
					commentLine := pass.Fset.Position(c.Pos()).Line
					// Check if comment is on the same line as Elem (inline comment)
					if commentLine == elemLine {
						hasComment = true
						break
					}
				}
				if hasComment {
					break
				}
			}

			if !hasComment {
				pos := pass.Fset.Position(schemaLit.Pos())
				if loader.ShouldReport(pos.Filename, pos.Line) {
					if propertyName := cached.PropertyName; propertyName != "" {
						pass.Reportf(schemaLit.Pos(), "%s: field `%s` has %s with only one nested property - consider %s or add inline comment explaining why (e.g., %s)\n",
							azsd001Name, propertyName,
							helper.IssueLine("MaxItems: 1"),
							helper.FixedCode("flattening"),
							helper.FixedCode("'// Additional properties will be added per service team confirmation'"))
					} else {
						pass.Reportf(schemaLit.Pos(), "%s: field has %s with only one nested property - consider %s or add inline comment explaining why (e.g., %s)\n",
							azsd001Name,
							helper.IssueLine("MaxItems: 1"),
							helper.FixedCode("flattening"),
							helper.FixedCode("'// Additional properties will be added per service team confirmation'"))
					}
				}
			}
		}
	}

	return nil, nil
}
