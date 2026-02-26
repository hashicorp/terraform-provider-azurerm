// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package passes

import (
	"go/ast"
	"sort"
	"strings"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/loader"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/passes/schema"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const AZNR001Doc = `check for top-level Schema field ordering

The AZNR001 analyzer reports cases of schemas where fields are not ordered correctly.

When git filter is applied, it only works on newly created files.

Schema fields should be ordered as follows:

1. Required fields in their original order
2. 'resource_group_name' must come before 'location' if both are required
3. Optional fields, sorted alphabetically (unless they appear before required fields)
4. Computed fields, sorted alphabetically (with 'location' first if computed)
5. 'tags' field must be at the end

Special cases:
- Schemas with 'name' field which is optional are skipped
- If optional fields appear before required fields, their original order is preserved
  (This happens when some resources have optional fields as part of the resource ID components)
- The expected order assumes ID fields are in the correct order; ID field ordering is not validated
- Nested schemas are not validated by this rule`

const aznr001Name = "AZNR001"

const (
	fieldName              = "name"
	fieldResourceGroupName = "resource_group_name"
	fieldLocation          = "location"
	fieldTags              = "tags"
)

var AZNR001Analyzer = &analysis.Analyzer{
	Name: aznr001Name,
	Doc:  AZNR001Doc,
	Run:  runAZNR001,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		schema.CompleteSchemaAnalyzer,
		commentignore.Analyzer,
	},
}

func runAZNR001(pass *analysis.Pass) (interface{}, error) {
	ignorer, ok := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	if !ok {
		return nil, nil
	}
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}
	completeSchemaInfo, ok := pass.ResultOf[schema.CompleteSchemaAnalyzer].(*schema.CompleteSchemaInfo)
	if !ok {
		return nil, nil
	}

	nodeFilter := []ast.Node{(*ast.CompositeLit)(nil)}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return
		}

		if ignorer.ShouldIgnore(aznr001Name, comp) {
			return
		}

		// Apply filename filtering
		filename := pass.Fset.Position(comp.Pos()).Filename
		if !loader.IsNewFile(filename) || !helper.IsResourceOrDataSourceFile(filename) {
			return
		}

		// Check if it's a schema map
		if !helper.IsSchemaMap(comp, pass.TypesInfo) {
			return
		}

		// Extract schema fields
		fields := completeSchemaInfo.SchemaFields[comp.Pos()]
		if len(fields) == 0 {
			return
		}

		// Check if this schema is nested within an Elem field
		isNested := false
		for _, f := range pass.Files {
			fPos := pass.Fset.Position(f.Pos())
			if fPos.Filename == filename {
				isNested = helper.IsNestedSchemaMap(f, comp)
				break
			}
		}
		// Skip nested schema for now, as this rule seems to apply for the top-level only
		// The nested schema order is not consistent
		if isNested {
			return
		}

		// Check for ordering issues
		expectedOrder, issue := checkAZNR001OrderingIssues(fields)
		if issue != "" {
			actualOrder := make([]string, len(fields))
			for i, f := range fields {
				actualOrder[i] = f.Name
			}
			pass.Reportf(comp.Pos(), "%s: %s\nExpected order (assuming ID fields are correct):\n  %s\nActual order:\n  %s\n",
				aznr001Name, issue,
				helper.FixedCode(strings.Join(expectedOrder, ", ")),
				helper.IssueLine(strings.Join(actualOrder, ", ")))
		}
	})

	return nil, nil
}

func checkAZNR001OrderingIssues(fields []helper.SchemaFieldInfo) ([]string, string) {
	if len(fields) == 0 {
		return nil, ""
	}

	// Skip if name field is optional + ExactlyOneOf (e.g., name/name_regex pattern (image_data_source.go))
	if hasOptionalName(fields) {
		return nil, ""
	}

	expectedOrder := getAZNR001ExpectedOrder(fields)
	return expectedOrder, validateAZNR001Order(fields, expectedOrder)
}

func getAZNR001ExpectedOrder(fields []helper.SchemaFieldInfo) []string {
	locationIsComputed := isLocationComputed(fields)
	hasOptionalBeforeRequired := checkOptionalBeforeRequired(fields)

	var result []string
	var computedFields []string
	var tagsField string

	if hasOptionalBeforeRequired {
		// Preserve original order for required and optional fields
		for _, field := range fields {
			if field.Name == fieldTags {
				tagsField = field.Name
				continue
			}
			if field.Name == fieldLocation && locationIsComputed {
				continue
			}
			if field.SchemaInfo != nil {
				schema := field.SchemaInfo.Schema
				// Only computed-only fields are separated
				if schema.Computed && !schema.Optional && !schema.Required {
					computedFields = append(computedFields, field.Name)
				} else {
					result = append(result, field.Name)
				}
			}
		}
	} else {
		// Normal case: categorize and sort
		requiredFields, optionalFields, computed, tags := categorizeFields(fields, locationIsComputed)
		computedFields = computed
		tagsField = tags

		// Fix resource_group_name < location order if both are required
		fixRequiredFieldOrder(requiredFields)

		result = append(result, requiredFields...)
		sort.Strings(optionalFields)
		result = append(result, optionalFields...)
	}

	// Add computed fields (sorted, with location first if computed)
	result = appendComputedFields(result, computedFields, locationIsComputed)

	if tagsField != "" {
		result = append(result, tagsField)
	}

	return result
}

func isLocationComputed(fields []helper.SchemaFieldInfo) bool {
	for _, field := range fields {
		if field.Name == fieldLocation && field.SchemaInfo != nil {
			schema := field.SchemaInfo.Schema
			// Only computed-only location
			return schema.Computed && !schema.Optional && !schema.Required
		}
	}
	return false
}

func checkOptionalBeforeRequired(fields []helper.SchemaFieldInfo) bool {
	lastOptionalIdx := -1
	firstRequiredIdx := -1
	for i, field := range fields {
		if field.SchemaInfo != nil {
			if field.SchemaInfo.Schema.Optional && lastOptionalIdx == -1 {
				lastOptionalIdx = i
			}
			if field.SchemaInfo.Schema.Required && firstRequiredIdx == -1 {
				firstRequiredIdx = i
			}
		}
	}
	return lastOptionalIdx != -1 && firstRequiredIdx != -1 && lastOptionalIdx < firstRequiredIdx
}

func categorizeFields(fields []helper.SchemaFieldInfo, locationIsComputed bool) (required, optional, computed []string, tags string) {
	for _, field := range fields {
		// Handle tags field separately
		if field.Name == fieldTags {
			tags = field.Name
			continue
		}

		// Skip location if it's computed-only (will be added at the beginning of computed fields)
		if field.Name == fieldLocation && locationIsComputed {
			continue
		}

		if field.SchemaInfo != nil {
			schema := field.SchemaInfo.Schema
			switch {
			case schema.Required:
				required = append(required, field.Name)
			case schema.Optional:
				optional = append(optional, field.Name)
			case schema.Computed && !schema.Optional && !schema.Required:
				// Only computed-only fields go here
				computed = append(computed, field.Name)
			}
		}
	}
	return
}

func fixRequiredFieldOrder(requiredFields []string) {
	// Only fix resource_group_name < location order if both are required
	rgPos := -1
	locPos := -1
	for i, name := range requiredFields {
		if name == fieldResourceGroupName {
			rgPos = i
		}
		if name == fieldLocation {
			locPos = i
		}
	}
	// Swap if location comes before resource_group_name
	if rgPos != -1 && locPos != -1 && locPos < rgPos {
		requiredFields[rgPos], requiredFields[locPos] = requiredFields[locPos], requiredFields[rgPos]
	}
}

func appendComputedFields(result, computedFields []string, locationIsComputed bool) []string {
	// Add location at the beginning of computed fields if it's computed
	if locationIsComputed {
		result = append(result, fieldLocation)
	}
	sort.Strings(computedFields)
	result = append(result, computedFields...)
	return result
}

func validateAZNR001Order(fields []helper.SchemaFieldInfo, expectedOrder []string) string {
	if len(fields) != len(expectedOrder) {
		// Skip if len is not equal, it happens when it's failed to extract field's properties;
		// it might because the schema is defined in another package, except commonschema
		return ""
	}

	// Compare actual order with expected order
	actualOrder := make([]string, len(fields))
	for i, field := range fields {
		actualOrder[i] = field.Name
	}

	// Check if actual order matches expected order
	for i := 0; i < len(actualOrder); i++ {
		if actualOrder[i] != expectedOrder[i] {
			return "schema fields are not in the expected order, please double check the order as mentioned in guide-new-resource.md or guide-new-data-source.md"
		}
	}

	return ""
}

func hasOptionalName(fields []helper.SchemaFieldInfo) bool {
	for _, field := range fields {
		if field.Name == fieldName && field.SchemaInfo != nil {
			// Check if it's optional
			if field.SchemaInfo.Schema.Optional {
				return true
			}
		}
	}
	return false
}
