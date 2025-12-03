// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S011 validates and fixes possible values in documentation match schema
type S011 struct{}

var _ Rule = new(S011)

func (s S011) ID() string {
	return "S011"
}

func (s S011) Name() string {
	return "Arguments Possible Values Consistency"
}

func (s S011) Description() string {
	return "Validates that possible values in documentation match schema definition"
}

func (s S011) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	if d.Type == data.ResourceTypeData {
		return nil
	}

	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return nil
	}

	return forEachSchemaProperty(d, "", d.SchemaProperties, d.DocumentArguments, d.DocumentArguments.BlockDefinitions,
		func(fullPath string, schemaProperty *models.SchemaProperty, docProperty *models.DocumentProperty) error {
			return s.checkPropertyPossibleValues(d, fullPath, schemaProperty, docProperty, fix)
		},
	)
}

// checkPropertyPossibleValues compares possible values between schema and documentation
func (s S011) checkPropertyPossibleValues(
	d *data.TerraformNodeData,
	path string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
	fix bool,
) error {
	schemaPossibleValues := schemaProperty.PossibleValues
	docPossibleValues := docProperty.PossibleValues()

	// Skip if schema has no possible values
	// This allows documentation to provide more information than schema
	if len(schemaPossibleValues) == 0 {
		return nil
	}

	// Skip if document property should be skipped (multiple possible value sections)
	if docProperty.ShouldSkip() {
		return nil
	}

	missed, spare := sliceDiff(schemaPossibleValues, docPossibleValues)

	if len(missed) == 0 && len(spare) == 0 {
		return nil
	}

	// Check if values may exist in doc but not in code format
	if mayExistsInDoc(docProperty.Content, schemaPossibleValues) {
		return nil
	}

	// Skip if this property has version changes in upgrade guide
	// This is because it's not consistent whether possible values change should be documented
	// when it's in the upgrade feature flag
	if HasVersionChanges(d.Document.Path, d.Name, path) {
		return nil
	}

	origLine := docProperty.Content
	var message string
	var fixedLine string

	if len(missed) > 0 {
		message = fmt.Sprintf("`%s` missing possible values in documentation: %s", util.Bold(path), formatPossibleValues(missed))
	}
	if len(spare) > 0 {
		if message != "" {
			message += "; "
		}
		message += fmt.Sprintf("`%s` has extra possible values in documentation: %s", util.Bold(path), formatPossibleValues(spare))
	}

	fixedLine = s.applyPossibleValuesFix(origLine, docProperty, schemaPossibleValues)
	if fix {
		s.applyFix(d, docProperty, schemaPossibleValues)
	}

	return NewValidationIssue(
		s.ID(),
		s.Name(),
		path,
		message,
		d.Document.Path,
		origLine,
		fixedLine,
	)
}

// sliceDiff compares two slices and returns values that are in want but not in got (missed),
// and values that are in got but not in want (spare)
func sliceDiff(want, got []string) (missed, spare []string) {
	if len(want) == 0 {
		return
	}

	// Convert to lowercase for case-insensitive comparison
	wantMap := make(map[string]string)
	for _, v := range want {
		wantMap[strings.ToLower(v)] = v
	}

	gotMap := make(map[string]string)
	for _, v := range got {
		gotMap[strings.ToLower(v)] = v
	}

	// Find missed values (in want but not in got)
	for lower, original := range wantMap {
		if _, ok := gotMap[lower]; !ok {
			missed = append(missed, original)
		}
	}

	// Find spare values (in got but not in want)
	for lower, original := range gotMap {
		if _, ok := wantMap[lower]; !ok {
			spare = append(spare, original)
		}
	}

	return
}

// mayExistsInDoc returns true if all values exist in the doc line (even if not in code format)
func mayExistsInDoc(docLine string, want []string) bool {
	for _, val := range want {
		if !strings.Contains(docLine, val) {
			return false
		}
	}
	return true
}

// formatPossibleValues formats a slice of values for display in error messages
func formatPossibleValues(values []string) string {
	codes := make([]string, 0, len(values))
	for _, val := range values {
		codes = append(codes, "`"+val+"`")
	}
	return "[" + strings.Join(codes, ", ") + "]"
}

// applyPossibleValuesFix generates the fixed line with updated possible values
func (s S011) applyPossibleValuesFix(line string, docProperty *models.DocumentProperty, newPossibleValues []string) string {
	if len(newPossibleValues) == 0 {
		return line
	}

	sorted := sortPossibleValues(newPossibleValues)
	formattedValues := formatPossibleValuesInline(sorted)
	var newLine string

	// Check if doc already has possible values defined
	if len(docProperty.PossibleValues()) > 0 && docProperty.EnumStart > 0 {
		// Replace existing possible values
		newLine = line[:docProperty.EnumStart]
		if len(sorted) == 1 {
			newLine += "The only possible value is "
		} else {
			newLine += "Possible values are "
		}
		newLine += formattedValues

		if docProperty.EnumEnd > 0 && docProperty.EnumEnd < len(line) {
			newLine += line[docProperty.EnumEnd:]
		} else {
			newLine += "."
		}
	} else {
		// Add possible values to the line
		// Find a good insertion point (before "Defaults to" or "Changing this forces")
		idx := strings.Index(line, "Defaults to")
		if idx < 0 {
			idx = strings.Index(line, "Changing this forces")
		}

		if idx > 0 {
			// Insert before "Defaults to" or "Changing this forces"
			newLine = strings.TrimRight(line[:idx], " ")
			newLine += " "
		} else {
			// No specific insertion point found, append at the end
			// Preserve the original line but ensure it ends properly
			newLine = strings.TrimRight(line, "\n")
			// Ensure there's a sentence ending before adding new content
			if !strings.HasSuffix(newLine, ".") && !strings.HasSuffix(newLine, "!") && !strings.HasSuffix(newLine, "?") {
				newLine += "."
			}
			newLine += " "
		}

		if len(sorted) == 1 {
			newLine += "The only possible value is "
		} else {
			newLine += "Possible values are "
		}
		newLine += formattedValues
		newLine += "."

		if idx > 0 {
			newLine += " " + line[idx:]
		}
	}

	return newLine
}

// formatPossibleValuesInline formats values for inline use in documentation
func formatPossibleValuesInline(values []string) string {
	if len(values) == 0 {
		return ""
	}

	sorted := sortPossibleValues(values)
	res := make([]string, len(sorted))
	for idx, val := range sorted {
		res[idx] = "`" + val + "`"
	}

	if len(res) == 1 {
		return res[0]
	}

	s := res[0]
	if len(res) >= 3 {
		s = strings.Join(res[:len(res)-1], ", ")
	}
	s += " and " + res[len(res)-1]
	return s
}

// sortPossibleValues sorts possible values in a sensible order
// For weekdays, it maintains chronological order (Monday -> Sunday)
// For other values, it sorts alphabetically (case-insensitive)
func sortPossibleValues(values []string) []string {
	if len(values) <= 1 {
		return values
	}

	// Check if all values are weekdays
	weekdayOrder := map[string]int{
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sunday":    7,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	allWeekdays := true
	for _, v := range values {
		if _, ok := weekdayOrder[v]; !ok {
			allWeekdays = false
			break
		}
	}

	sorted := make([]string, len(values))
	copy(sorted, values)

	if allWeekdays {
		sort.Slice(sorted, func(i, j int) bool {
			return weekdayOrder[sorted[i]] < weekdayOrder[sorted[j]]
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i]) < strings.ToLower(sorted[j])
		})
	}

	return sorted
}

// applyFix applies the possible values fix to the document
func (s S011) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, newPossibleValues []string) {
	if d.Document == nil {
		return
	}

	argsSection := d.Document.GetArgumentsSection()
	if argsSection == nil {
		return
	}

	content := argsSection.GetContent()
	lineIdx := docProperty.Line
	if lineIdx >= 0 && lineIdx < len(content) {
		line := content[lineIdx]
		fixedLine := s.applyPossibleValuesFix(line, docProperty, newPossibleValues)
		content[lineIdx] = fixedLine
		argsSection.SetContent(content)
		d.Document.HasChange = true
	}
}
