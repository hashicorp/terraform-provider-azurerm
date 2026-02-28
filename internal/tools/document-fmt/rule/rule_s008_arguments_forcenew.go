// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S008 validates and fixes ForceNew markers in documentation match schema
type S008 struct{}

var _ Rule = new(S008)

var forceNewReg = regexp.MustCompile(` ?Changing.*forces? a [^.]*(\.|$)`)

func (s S008) ID() string {
	return "S008"
}

func (s S008) Name() string {
	return "Arguments ForceNew Consistency"
}

func (s S008) Description() string {
	return "Validates that ForceNew markers in documentation match schema definition"
}

func (s S008) Run(d *data.TerraformNodeData, fix bool) []error {
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
			return s.checkPropertyForceNew(d, fullPath, schemaProperty, docProperty, fix)
		},
	)
}

// checkPropertyForceNew checks and optionally fixes ForceNew marker for a single property
func (s S008) checkPropertyForceNew(
	d *data.TerraformNodeData,
	fullPath string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
	fix bool,
) error {
	if docProperty == nil || schemaProperty == nil {
		return nil
	}

	// Skip resource_group_name as per existing logic
	if lastPathSegment(fullPath) == "resource_group_name" {
		return nil
	}

	// Check: ForceNew markers
	if schemaProperty.ForceNew != docProperty.ForceNew {
		if schemaProperty.ForceNew && !docProperty.ForceNew {
			// Should add ForceNew marker
			origLine := strings.TrimRight(docProperty.Content, "\n")
			fixedLine := s.fixForceNew(origLine, true)
			issue := NewValidationIssue(
				s.ID(),
				s.Name(),
				fullPath,
				fmt.Sprintf("`%s` should be marked as ForceNew", util.Bold(fullPath)),
				d.Document.Path,
				origLine,
				fixedLine,
			)

			if fix {
				s.applyFix(d, docProperty, true)
			}
			return issue
		} else if docProperty.ForceNew && !schemaProperty.ForceNew {
			// Should remove ForceNew marker
			origLine := strings.TrimRight(docProperty.Content, "\n")
			fixedLine := s.fixForceNew(origLine, false)
			issue := NewValidationIssue(
				s.ID(),
				s.Name(),
				fullPath,
				fmt.Sprintf("`%s` should NOT be marked as ForceNew", util.Bold(fullPath)),
				d.Document.Path,
				origLine,
				fixedLine,
			)

			if fix {
				s.applyFix(d, docProperty, false)
			}
			return issue
		}
	}

	return nil
}

// applyFix applies ForceNew fix to the document
func (s S008) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, shouldAdd bool) {
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
		fixedLine := s.fixForceNew(line, shouldAdd)
		content[lineIdx] = fixedLine
		argsSection.SetContent(content)
		d.Document.HasChange = true
	}
}

func (s S008) fixForceNew(line string, shouldAdd bool) string {
	if shouldAdd {
		// Add ForceNew message if not present
		line = strings.TrimRight(line, " \t\r\n")
		if strings.HasSuffix(line, ",") {
			line = line[:len(line)-1] + "."
		} else if !strings.HasSuffix(line, ".") {
			line += "."
		}
		line += " Changing this forces a new resource to be created."
	} else {
		// Remove ForceNew message
		line = forceNewReg.ReplaceAllString(line, "")
	}
	return line
}

// lastPathSegment extracts the last segment of a dotted path
func lastPathSegment(path string) string {
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[idx+1:]
	}
	return path
}
