package passes

import (
	"go/ast"
	"strings"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	localschema "github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
)

const AZBP002Doc = `check Optional+Computed fields follow conventions

The AZBP002 analyzer checks that fields marked as both Optional and Computed:
1. Have properties in sequence: Optional, Comment, Computed
2. Have a comment starting with "// NOTE: O+C " explaining why

Example violation:
  "field": {
      Type:     schema.TypeString,
      Optional: true,
      Computed: true,  // Missing NOTE: O+C comment
  }

Valid usage:
  "field": {
      Type:     schema.TypeString,
      Optional: true,
      // NOTE: O+C - field can be set by user or computed from API when not provided
      Computed: true,
  }`

const azbp002Name = "AZBP002"

var AZBP002Analyzer = &analysis.Analyzer{
	Name:     azbp002Name,
	Doc:      AZBP002Doc,
	Run:      runAZBP002,
	Requires: []*analysis.Analyzer{localschema.LocalAnalyzer},
}

func runAZBP002(pass *analysis.Pass) (interface{}, error) {
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

		// Only check fields that are both Optional and Computed
		if !schemaInfo.Schema.Optional || !schemaInfo.Schema.Computed {
			continue
		}

		// Get positions of Optional and Computed fields from cached SchemaInfo
		optionalKV := schemaInfo.Fields[schema.SchemaFieldOptional]
		computedKV := schemaInfo.Fields[schema.SchemaFieldComputed]
		if optionalKV == nil || computedKV == nil {
			continue
		}

		optionalPos := optionalKV.Pos()
		computedPos := computedKV.Pos()

		// Check order: Optional should come before Computed
		if optionalPos > computedPos {
			pos := pass.Fset.Position(schemaLit.Pos())
			if loader.ShouldReport(pos.Filename, pos.Line) {
				pass.Reportf(schemaLit.Pos(), "%s: field has %s and %s in wrong order (%s must come before %s)\n",
					azbp002Name,
					helper.FixedCode("Optional"), helper.IssueLine("Computed"),
					helper.FixedCode("Optional"), helper.IssueLine("Computed"))
			}
			continue
		}

		// Check for NOTE: O+C comment between Optional and Computed
		filename := pass.Fset.Position(schemaLit.Pos()).Filename
		optionalLine := pass.Fset.Position(optionalPos).Line
		computedLine := pass.Fset.Position(computedPos).Line

		hasOCComment := false
		comments := fileCommentsMap[filename]
		for _, cg := range comments {
			for _, c := range cg.List {
				commentLine := pass.Fset.Position(c.Pos()).Line
				if commentLine > optionalLine && commentLine < computedLine {
					if strings.Contains(c.Text, "NOTE: O+C") {
						hasOCComment = true
						break
					}
				}
			}
			if hasOCComment {
				break
			}
		}

		if !hasOCComment {
			pos := pass.Fset.Position(schemaLit.Pos())
			if loader.ShouldReport(pos.Filename, pos.Line) {
				if propertyName := cached.PropertyName; propertyName != "" {
					pass.Reportf(schemaLit.Pos(), "%s: field `%s` is Optional+Computed but missing required comment. Add %s between Optional and Computed\n",
						azbp002Name, propertyName, helper.FixedCode("'// NOTE: O+C - <explanation>'"))
				} else {
					pass.Reportf(schemaLit.Pos(), "%s: field is Optional+Computed but missing required comment. Add %s between Optional and Computed\n",
						azbp002Name, helper.FixedCode("'// NOTE: O+C - <explanation>'"))
				}
			}
		}
	}

	return nil, nil
}
