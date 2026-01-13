package passes

import (
	"go/ast"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	localschema "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
)

const AZBP001Doc = `check that all String arguments have validation

The AZBP001 analyzer reports cases where String type schema fields 
(Required or Optional) do not have a ValidateFunc.

Example violations:
  "name": {
      Type:     schema.TypeString,
      Required: true,
      // Missing ValidateFunc!
  }
  
Valid usage:
  "name": {
      Type:         schema.TypeString,
      Required:     true,
      ValidateFunc: validation.StringIsNotEmpty,
  }
  
  "description": {
      Type:     schema.TypeString,
      Computed: true,  // OK - computed-only fields don't need validation
  }`

const azbp001Name = "AZBP001"

var AZBP001Analyzer = &analysis.Analyzer{
	Name:     azbp001Name,
	Doc:      AZBP001Doc,
	Run:      runAZBP001,
	Requires: []*analysis.Analyzer{localschema.LocalAnalyzer},
}

func runAZBP001(pass *analysis.Pass) (interface{}, error) {
	schemaInfoCache, ok := pass.ResultOf[localschema.LocalAnalyzer].(map[*ast.CompositeLit]*localschema.LocalSchemaInfoWithName)
	if !ok {
		return nil, nil
	}

	for schemaLit, cached := range schemaInfoCache {
		schemaInfo := cached.Info

		// Type check: only check String fields
		if !schemaInfo.IsType(schema.SchemaValueTypeString) {
			continue
		}

		// Skip computed-only fields (no Required or Optional)
		if schemaInfo.Schema.Computed && !schemaInfo.Schema.Required && !schemaInfo.Schema.Optional {
			continue
		}

		// Check if validation exists
		hasValidation := schemaInfo.DeclaresField(schema.SchemaFieldValidateFunc)

		if !hasValidation {
			pos := pass.Fset.Position(schemaLit.Pos())
			if loader.ShouldReport(pos.Filename, pos.Line) {
				if propertyName := cached.PropertyName; propertyName != "" {
					pass.Reportf(schemaLit.Pos(), "%s: string argument `%s` %s\n",
						azbp001Name, propertyName, helper.FixedCode("must have ValidateFunc"))
				} else {
					pass.Reportf(schemaLit.Pos(), "%s: string argument %s\n",
						azbp001Name, helper.FixedCode("must have ValidateFunc"))
				}
			}
		}
	}

	return nil, nil
}
