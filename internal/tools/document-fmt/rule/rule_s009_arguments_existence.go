// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/mdparser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S009 validates that properties exist in both schema and documentation
type S009 struct{}

var _ Rule = new(S009)

func (s S009) ID() string {
	return "S009"
}

func (s S009) Name() string {
	return "Arguments Exist in Schema and Document"
}

func (s S009) Description() string {
	return "Validates that all arguments exist in both schema and documentation, checks for missing properties and potential misspellings"
}

func (s S009) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	if d.Type == data.ResourceTypeData {
		return nil
	}

	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return nil
	}

	var errs []error
	resourceType := d.Name

	// First pass: check schema properties missing in documentation
	errs = append(errs, s.checkMissingInDoc(d, "", d.SchemaProperties, d.DocumentArguments, d.DocumentArguments.BlockDefinitions, resourceType)...)

	// Second pass: check documentation properties missing in schema (potential misspellings)
	errs = append(errs, s.checkMissingInSchema(d, "", d.DocumentArguments, d.SchemaProperties, resourceType)...)

	return errs
}

// checkMissingInDoc checks if schema properties are documented and validates block declarations
func (s S009) checkMissingInDoc(
	d *data.TerraformNodeData,
	parentPath string,
	schema *models.SchemaProperties,
	documentation *models.DocumentProperties,
	blockDefinitions map[string]*models.DocumentProperty,
	resourceType string,
) []error {
	var errs []error

	if schema == nil {
		return errs
	}

	for name, schemaProperty := range schema.Objects {
		// Skip computed properties and `id` field
		if schemaProperty.Computed {
			continue
		}
		if name == "id" {
			continue
		}
		// Skip deprecated properties
		if schemaProperty.Deprecated {
			continue
		}

		fullPath := name
		if parentPath != "" {
			fullPath = parentPath + "." + name
		}

		// Skip properties in skip config
		if SkipProp(s.ID(), resourceType, fullPath) {
			continue
		}

		// Check if property exists in documentation
		docProperty := documentation.Objects[name]
		if docProperty == nil {
			issue := NewValidationIssue(
				s.ID(),
				s.Name(),
				fullPath,
				fmt.Sprintf("`%s` exists in schema but is MISSING from documentation", util.Bold(fullPath)),
				d.Document.Path,
				"",
				"",
			)
			errs = append(errs, issue)
			continue
		}

		// Check for block type declarations (nested properties)
		if schemaProperty.Nested != nil && len(schemaProperty.Nested.Objects) > 0 {
			// Check if the field is marked as a block in documentation
			if !docProperty.Block {
				issue := NewValidationIssue(
					s.ID(),
					s.Name(),
					fullPath,
					fmt.Sprintf("`%s` should be declared as a BLOCK", util.Bold(fullPath)),
					d.Document.Path,
					"",
					"",
				)
				errs = append(errs, issue)
				continue
			}

			if docProperty.Nested == nil || len(docProperty.Nested.Objects) == 0 {
				// For some blocks sharing same sub-fields, they are defined in a shared block section
				if docProperty.BlockTypeName != docProperty.Name {
					linkedDocProperty := blockDefinitions[docProperty.BlockTypeName]
					if linkedDocProperty != nil && linkedDocProperty.Nested != nil && len(linkedDocProperty.Nested.Objects) > 0 {
						// Recursively check nested properties in shared block
						errs = append(errs, s.checkMissingInDoc(d, fullPath, schemaProperty.Nested, linkedDocProperty.Nested, blockDefinitions, resourceType)...)
						continue
					}
				}

				issue := NewValidationIssue(
					s.ID(),
					s.Name(),
					fullPath,
					fmt.Sprintf("`%s` block is MISSING nested properties", util.Bold(fullPath)),
					d.Document.Path,
					"",
					"",
				)
				errs = append(errs, issue)
				continue
			}

			// Recursively check nested properties
			if docProperty.Nested != nil {
				errs = append(errs, s.checkMissingInDoc(d, fullPath, schemaProperty.Nested, docProperty.Nested, blockDefinitions, resourceType)...)
			}
		}
	}

	return errs
}

// checkMissingInSchema checks if documented properties exist in schema
func (s S009) checkMissingInSchema(
	d *data.TerraformNodeData,
	parentPath string,
	documentation *models.DocumentProperties,
	schema *models.SchemaProperties,
	resourceType string,
) []error {
	var errs []error

	if documentation == nil {
		return errs
	}

	for name, docProperty := range documentation.Objects {
		if name == "id" {
			continue
		}

		fullPath := name
		if parentPath != "" {
			fullPath = parentPath + "." + name
		}

		// Skip properties in skip config
		if SkipProp(s.ID(), resourceType, fullPath) {
			continue
		}

		if len(docProperty.ParseErrors) > 0 {
			for _, parseErr := range docProperty.ParseErrors {
				if strings.Contains(parseErr, mdparser.MisspelNameOfProperty) {
					issue := NewValidationIssue(
						s.ID(),
						s.Name(),
						fullPath,
						fmt.Sprintf("`%s` does NOT exist in schema - possible misspelling?", util.Bold(fullPath)),
						d.Document.Path,
						"",
						"",
					)
					errs = append(errs, issue)
				}
			}
			continue
		}

		if strings.Contains(strings.ToLower(docProperty.Content), "deprecated") {
			continue
		}

		// Check if property exists in schema
		var schemaProperty *models.SchemaProperty
		if schema != nil {
			schemaProperty = schema.Objects[name]
		}

		if schemaProperty == nil {
			issue := NewValidationIssue(
				s.ID(),
				s.Name(),
				fullPath,
				fmt.Sprintf("`%s` exists in documentation but NOT in schema", util.Bold(fullPath)),
				d.Document.Path,
				"",
				"",
			)
			errs = append(errs, issue)
			continue
		}

		// Check if document marks field as block but schema doesn't have nested properties
		if docProperty.Block && (schemaProperty.Nested == nil || len(schemaProperty.Nested.Objects) == 0) {
			issue := NewValidationIssue(
				s.ID(),
				s.Name(),
				fullPath,
				fmt.Sprintf("The document incorrectly implies `%s` is a BLOCK (contains phrases like 'as defined below')", util.Bold(fullPath)),
				d.Document.Path,
				"",
				"",
			)
			errs = append(errs, issue)
			continue
		}

		// Recursively check nested properties
		if docProperty.Nested != nil && len(docProperty.Nested.Objects) > 0 {
			if schemaProperty.Nested != nil {
				errs = append(errs, s.checkMissingInSchema(d, fullPath, docProperty.Nested, schemaProperty.Nested, resourceType)...)
			}
		}
	}

	return errs
}
