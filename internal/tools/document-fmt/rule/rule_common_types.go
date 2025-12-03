// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// ValidationIssue represents a validation issue with before/after comparison
type ValidationIssue struct {
	RuleID      string
	RuleName    string
	PropertyKey string
	Message     string
	FileName    string
	OrigLine    string // Original line content
	FixedLine   string // Fixed line content (empty if no fix available)
}

func (vi *ValidationIssue) Error() string {
	var result strings.Builder

	// Format: "S007: fileName: message"
	result.WriteString(vi.RuleID)
	result.WriteString(": ")
	result.WriteString(getRelevantPath(vi.FileName))
	result.WriteString(": ")
	result.WriteString(vi.Message)

	// If we have both original and fixed lines, show the comparison
	if vi.OrigLine != "" && vi.FixedLine != "" && vi.OrigLine != vi.FixedLine {
		result.WriteString("\n     ")
		result.WriteString(util.Yellow(strings.TrimRight(vi.OrigLine, "\n")))
		result.WriteString("\n  => ")
		result.WriteString(util.Green(strings.TrimRight(vi.FixedLine, "\n")))
		result.WriteString("\n")
	}

	return result.String()
}

// NewValidationIssue creates a new validation issue
func NewValidationIssue(ruleID, ruleName, propertyKey, message, fileName, origLine, fixedLine string) *ValidationIssue {
	return &ValidationIssue{
		RuleID:      ruleID,
		RuleName:    ruleName,
		PropertyKey: propertyKey,
		Message:     message,
		FileName:    fileName,
		OrigLine:    origLine,
		FixedLine:   fixedLine,
	}
}

// getRelevantPath extracts the relevant part of file path for display
func getRelevantPath(fullPath string) string {
	path := filepath.ToSlash(fullPath)

	if idx := strings.Index(path, "website/docs/"); idx >= 0 {
		return path[idx:]
	}

	return filepath.Base(fullPath)
}

// PropertyValidatorFunc validates a single property and returns an error if validation fails
type PropertyValidatorFunc func(
	fullPath string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
) error

// forEachSchemaProperty walks through schema properties and calls the validator function for each valid property.
func forEachSchemaProperty(
	ruleId string,
	d *data.TerraformNodeData,
	parentPath string,
	schema *models.SchemaProperties,
	documentation *models.DocumentProperties,
	blockDefinitions map[string]*models.DocumentProperty,
	validator PropertyValidatorFunc,
) []error {
	var errs []error

	if schema == nil || documentation == nil {
		return errs
	}

	for name, schemaProperty := range schema.Objects {
		// Common filtering: skip computed, id, deprecated
		if schemaProperty.Computed || name == "id" || schemaProperty.Deprecated {
			continue
		}

		fullPath := name
		if parentPath != "" {
			fullPath = parentPath + "." + name
		}

		// Skip properties in skip config
		if SkipProp(ruleId, d.Name, fullPath) {
			continue
		}

		docProperty := documentation.Objects[name]
		if docProperty == nil {
			continue
		}

		// Skip properties with parse errors
		if len(docProperty.ParseErrors) > 0 {
			continue
		}

		// Handle nested properties (blocks)
		if schemaProperty.Nested != nil && len(schemaProperty.Nested.Objects) > 0 {
			if !docProperty.Block {
				continue // Block declaration error, skip
			}

			// Get nested documentation (handle shared block definitions)
			nestedDocs := docProperty.Nested
			if docProperty.Nested == nil || len(docProperty.Nested.Objects) == 0 {
				if docProperty.BlockTypeName != docProperty.Name {
					if linkedDocProperty := blockDefinitions[docProperty.BlockTypeName]; linkedDocProperty != nil {
						nestedDocs = linkedDocProperty.Nested
					}
				}
			}

			// Recursively check nested properties
			if nestedDocs != nil && len(nestedDocs.Objects) > 0 {
				nestedErrs := forEachSchemaProperty(ruleId, d, fullPath, schemaProperty.Nested, nestedDocs, blockDefinitions, validator)
				errs = append(errs, nestedErrs...)
			}

			// Also validate the block property itself (for ForceNew checks, etc.)
			if err := validator(fullPath, schemaProperty, docProperty); err != nil {
				errs = append(errs, err)
			}
		} else {
			// For non-nested properties: validate directly
			if err := validator(fullPath, schemaProperty, docProperty); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

// DocPropertyValidatorFunc validates a single documented property and returns an error if validation fails
type DocPropertyValidatorFunc func(
	fullPath string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
) error

// forEachDocProperty walks through documented properties and calls the validator function for each.
func forEachDocProperty(
	ruleId string,
	d *data.TerraformNodeData,
	parentPath string,
	documentation *models.DocumentProperties,
	schema *models.SchemaProperties,
	validator DocPropertyValidatorFunc,
) []error {
	var errs []error

	if documentation == nil {
		return errs
	}

	for name, docProperty := range documentation.Objects {
		// Skip properties that are not marked as Required or Optional
		if !docProperty.Required && !docProperty.Optional {
			continue
		}
		if name == "id" {
			continue
		}

		fullPath := name
		if parentPath != "" {
			fullPath = parentPath + "." + name
		}

		// Skip properties in skip config
		if SkipProp(ruleId, d.Name, fullPath) {
			continue
		}

		var schemaProperty *models.SchemaProperty
		if schema != nil {
			schemaProperty = schema.Objects[name]
		}

		// Validate the property
		if err := validator(fullPath, schemaProperty, docProperty); err != nil {
			errs = append(errs, err)
		}

		// Recursively check nested properties
		if docProperty.Nested != nil && len(docProperty.Nested.Objects) > 0 {
			var nestedSchema *models.SchemaProperties
			if schemaProperty != nil {
				nestedSchema = schemaProperty.Nested
			}
			nestedErrs := forEachDocProperty(ruleId, d, fullPath, docProperty.Nested, nestedSchema, validator)
			errs = append(errs, nestedErrs...)
		}
	}

	return errs
}
