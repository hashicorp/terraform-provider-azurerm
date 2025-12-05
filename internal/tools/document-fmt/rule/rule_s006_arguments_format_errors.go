// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/mdparser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
)

// S006 validates and fixes Arguments format errors
type S006 struct{}

var _ Rule = new(S006)

var regIncorrectBlockMark = regexp.MustCompile(`(?: blocks?)? as (?:detailed|defined) (below|above)`)

func (s S006) ID() string {
	return "S006"
}

func (s S006) Name() string {
	return "Arguments Format Errors"
}

func (s S006) Description() string {
	return "Validates and fixes document formatting errors in Arguments section"
}

func (s S006) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	if d.Type == data.ResourceTypeData {
		return nil
	}

	if d.DocumentArguments == nil {
		return nil
	}

	return forEachDocProperty(s.ID(), d, "", d.DocumentArguments, d.SchemaProperties,
		func(fullPath string, schemaProperty *models.SchemaProperty, docProperty *models.DocumentProperty) error {
			return s.checkPropertyFormat(d, fullPath, schemaProperty, docProperty, fix)
		},
	)
}

func (s S006) checkPropertyFormat(
	d *data.TerraformNodeData,
	fullPath string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
	fix bool,
) error {
	// Handle parse errors
	if len(docProperty.ParseErrors) == 0 {
		return nil
	}

	// Skip Computed properties for now
	if schemaProperty.Computed {
		return nil
	}

	for _, parseErr := range docProperty.ParseErrors {
		// Check if "block is not defined" error should be converted to "incorrectly block marked"
		if strings.Contains(parseErr, mdparser.BlcokNotDefined) && schemaProperty != nil {
			// If schema property exists but is not a block, update the error type
			if schemaProperty.Nested == nil || len(schemaProperty.Nested.Objects) == 0 {
				parseErr = "incorrectly block marked"
			}
		}

		if strings.Contains(parseErr, "misspell of name from") {
			continue
		}

		var message string
		origLine := strings.TrimRight(docProperty.Content, "\n")
		fixedLine := s.getFixedLine(docProperty, parseErr)

		switch {
		case strings.Contains(parseErr, mdparser.IncorrectlyBlockMarked):
			message = fmt.Sprintf("The document incorrectly implies `%s` is a block (contains phrases like 'as defined below')", fullPath)
		case strings.Contains(parseErr, "duplicate"):
			message = fmt.Sprintf("%s: `%s`", parseErr, fullPath)
		case strings.Contains(parseErr, "no field name found"):
			message = fmt.Sprintf("following should be formatted as: `* `field` - (Required/Optional) Xxx...`\n  %s\n", docProperty.Content)
		default:
			message = parseErr
		}

		issue := NewValidationIssue(
			s.ID(),
			s.Name(),
			fullPath,
			message,
			d.Document.Path,
			origLine,
			fixedLine,
		)

		if fix {
			s.applyFix(d, docProperty, parseErr)
		}

		return issue
	}

	return nil
}

// applyFix applies format fix to the document
func (s S006) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, parseErr string) {
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
		fixedLine := s.fixFormatError(line, parseErr)

		// Note: We don't delete lines to avoid index invalidation for subsequent rules
		// Empty lines will be preserved to maintain line numbers
		if fixedLine != line {
			content[lineIdx] = fixedLine
			argsSection.SetContent(content)
			d.Document.HasChange = true
		}
	}
}

// getFixedLine returns the fixed version of the line for display purposes
func (s S006) getFixedLine(docProperty *models.DocumentProperty, parseErr string) string {
	if docProperty == nil {
		return ""
	}
	line := strings.TrimRight(docProperty.Content, "\n")
	return s.fixFormatError(line, parseErr)
}

// fixFormatError applies formatting fixes to a line based on the parse error
func (s S006) fixFormatError(line string, parseErr string) string {
	// Remove misleading star mark from Note lines
	if strings.HasPrefix(line, "* ~>") {
		return strings.TrimPrefix(line, "* ")
	}

	// Mark empty list markers as empty (but preserve line to avoid index issues)
	if strings.TrimSpace(line) == "*" {
		return "" // Return empty string, line will be cleared but not removed
	}

	// Remove incorrect block markers
	if strings.Contains(parseErr, mdparser.IncorrectlyBlockMarked) {
		return regIncorrectBlockMark.ReplaceAllLiteralString(line, "")
	}

	return line
}
