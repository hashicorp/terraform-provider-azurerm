// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/mdparser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S010 validates and fixes default values in documentation match schema
type S010 struct{}

var _ Rule = new(S010)

func (s S010) ID() string {
	return "S010"
}

func (s S010) Name() string {
	return "Arguments Default Value Consistency"
}

func (s S010) Description() string {
	return "Validates that default values in documentation match schema definition"
}

func (s S010) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	if d.Type == data.ResourceTypeData {
		return nil
	}

	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return nil
	}

	return forEachSchemaProperty(s.ID(), d, "", d.SchemaProperties, d.DocumentArguments, d.DocumentArguments.BlockDefinitions,
		func(fullPath string, schemaProperty *models.SchemaProperty, docProperty *models.DocumentProperty) error {
			return s.checkPropertyDefaultValue(d, fullPath, schemaProperty, docProperty, fix)
		},
	)
}

// checkPropertyDefaultValue compares default value between schema and documentation
func (s S010) checkPropertyDefaultValue(
	d *data.TerraformNodeData,
	path string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
	fix bool,
) error {
	var schemaDefault string
	if schemaProperty.DefaultValue != nil {
		schemaDefault = fmt.Sprintf("%v", schemaProperty.DefaultValue)
		if str, ok := schemaProperty.DefaultValue.(string); ok && str == "" {
			schemaDefault = `""` // empty string in code
		}
	}

	docDefault := docProperty.DefaultValue

	if s.defaultValuesMatch(schemaDefault, docDefault) {
		return nil
	}

	origLine := docProperty.Content
	var message string
	var fixedLine string

	// Case 1: Schema has default but doc doesn't, or values don't match
	if schemaDefault != "" {
		fixedLine = s.fixDefaultValue(origLine, schemaDefault)
		message = fmt.Sprintf("`%s` missing / wrong Default Value declaration", util.Bold(path))

		if fix {
			s.applyFix(d, docProperty, schemaDefault)
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

	// Case 2: Schema has no default but doc has one (and schema is not computed)
	if docDefault != "" && !schemaProperty.Computed {
		// Special case: boolean type with false default is often omitted
		if schemaProperty.Type == "Bool" && docDefault == "false" {
			return nil
		}

		fixedLine = s.fixDefaultValue(origLine, "")
		message = fmt.Sprintf("`%s` should NOT have default value", util.Bold(path))

		if fix {
			s.applyFix(d, docProperty, "")
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

	return nil
}

// defaultValuesMatch checks if two default values are equivalent
func (s S010) defaultValuesMatch(schemaDefault, docDefault string) bool {
	if schemaDefault == docDefault {
		return true
	}

	if schemaDefault == "" && docDefault == "" {
		return true
	}

	// Special case: false default for boolean types
	if schemaDefault == "false" && docDefault == "" {
		return true
	}

	// Try numeric comparison for numbers
	if schemaNum, e1 := strconv.ParseFloat(schemaDefault, 64); e1 == nil {
		if docNum, e2 := strconv.ParseFloat(docDefault, 64); e2 == nil {
			return int(schemaNum) == int(docNum)
		}
	}

	return false
}

// fixDefaultValue generates the fixed line with updated default value
func (s S010) fixDefaultValue(line string, newDefaultValue string) string {
	var newLine string

	if idxs := mdparser.DefaultsReg.FindStringSubmatchIndex(line); len(idxs) > 2 {
		if newDefaultValue == "" {
			newLine = line[:idxs[0]+1] + line[idxs[1]:]
		} else {
			newLine = line[:idxs[2]] + "`" + newDefaultValue + "`" + line[idxs[3]:]
		}
	} else if newDefaultValue != "" {
		newLine = strings.TrimSpace(line)
		if !strings.HasSuffix(newLine, ".") {
			newLine += "."
		}
		newLine += " Defaults to `" + newDefaultValue + "`."
	} else {
		// No default value to add or remove
		return line
	}

	return newLine
}

// applyFix applies the default value fix to the document
func (s S010) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, fixedLine string) {
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
		fixedLine := s.fixDefaultValue(line, fixedLine)
		content[lineIdx] = fixedLine
		argsSection.SetContent(content)
		d.Document.HasChange = true
	}
}
