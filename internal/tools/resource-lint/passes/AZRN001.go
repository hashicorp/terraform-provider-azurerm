package passes

import (
	"go/ast"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	localschema "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
)

const AZRN001Doc = `check that percentage properties use _percentage suffix instead of _in_percent

The AZRN001 analyzer reports when percentage property names use '_in_percent' 
suffix instead of the preferred '_percentage' suffix.

Example violations:
  "cpu_in_percent": {...}      // should be "cpu_percentage"
  "memory_in_percent": {...}   // should be "memory_percentage"

Valid usage:
  "cpu_percentage": {...}
  "memory_percentage": {...}`

const azrn001Name = "AZRN001"

var AZRN001Analyzer = &analysis.Analyzer{
	Name:     azrn001Name,
	Doc:      AZRN001Doc,
	Run:      runAZRN001,
	Requires: []*analysis.Analyzer{localschema.LocalAnalyzer},
}

func runAZRN001(pass *analysis.Pass) (interface{}, error) {
	schemaInfoCache, ok := pass.ResultOf[localschema.LocalAnalyzer].(map[*ast.CompositeLit]*localschema.LocalSchemaInfoWithName)
	if !ok {
		return nil, nil
	}

	for schemaLit, cached := range schemaInfoCache {
		fieldName := cached.PropertyName

		// Check if field name contains "_in_percent"
		if strings.Contains(fieldName, "_in_percent") {
			suggestedName := strings.ReplaceAll(fieldName, "_in_percent", "_percentage")
			pos := pass.Fset.Position(schemaLit.Pos())
			// Only report if this line is in the changed lines
			if loader.ShouldReport(pos.Filename, pos.Line) {
				pass.Reportf(schemaLit.Pos(), "%s: field %q should use %s suffix instead of %s (suggested: %q)\n",
					azrn001Name, fieldName,
					helper.FixedCode("'_percentage'"),
					helper.IssueLine("'_in_percent'"),
					suggestedName)
			}
		}
	}

	return nil, nil
}
