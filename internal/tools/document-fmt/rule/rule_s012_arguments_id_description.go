// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S012 validates that xxx_id fields have descriptions starting with "The ID of the"
type S012 struct{}

var _ Rule = new(S012)

func (s S012) ID() string   { return "S012" }
func (s S012) Name() string { return "Arguments ID Description Format" }
func (s S012) Description() string {
	return "Validates that xxx_id fields have descriptions starting with 'The ID of the'"
}

func (s S012) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists || d.Type == data.ResourceTypeData {
		return nil
	}
	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return nil
	}

	return forEachSchemaProperty(s.ID(), d, "", d.SchemaProperties, d.DocumentArguments, d.DocumentArguments.BlockDefinitions,
		func(fullPath string, schemaProperty *models.SchemaProperty, docProperty *models.DocumentProperty) error {
			return s.checkIdDescription(d, fullPath, docProperty, fix)
		},
	)
}

func (s S012) checkIdDescription(d *data.TerraformNodeData, fullPath string, docProperty *models.DocumentProperty, fix bool) error {
	if !strings.HasSuffix(strings.ToLower(docProperty.Name), "_id") {
		return nil
	}

	if !docProperty.Required && !docProperty.Optional {
		return nil
	}

	// Extract description after "(Required)" or "(Optional)"
	content := docProperty.Content
	marker, idx := s.findMarker(content)
	if idx == -1 {
		return nil
	}
	description := strings.TrimSpace(content[idx+len(marker):])

	if strings.HasPrefix(strings.ToLower(description), "the id of the") {
		return nil
	}

	// Generate fix suggestion
	fixedLine := s.fixIdDescription(docProperty.Name, content)
	if fix && fixedLine != "" {
		s.applyFix(d, docProperty, fixedLine)
	}

	return NewValidationIssue(s.ID(), s.Name(), fullPath,
		fmt.Sprintf("`%s` description should start with 'The ID of the ...'", util.Bold(fullPath)),
		d.Document.Path, content, fixedLine)
}

// findMarker finds "(Required) " or "(Optional) " in content and returns the marker and its index
func (s S012) findMarker(content string) (string, int) {
	if idx := strings.Index(content, "(Required) "); idx != -1 {
		return "(Required) ", idx
	}
	if idx := strings.Index(content, "(Optional) "); idx != -1 {
		return "(Optional) ", idx
	}
	return "", -1
}

// fixIdDescription generates a fixed description line with "The ID of the <ResourceType>." format
func (s S012) fixIdDescription(fieldName, content string) string {
	marker, idx := s.findMarker(content)
	if idx == -1 {
		return ""
	}

	prefix := content[:idx+len(marker)]
	description := content[idx+len(marker):]

	// Infer resource type from field name (e.g., "virtual_network_id" -> "virtual network")
	resourceType := strings.TrimSuffix(fieldName, "_id")
	resourceType = strings.ReplaceAll(resourceType, "_", " ")

	// Find the first sentence ending and preserve everything after it
	suffix := ""
	if dotIdx := strings.Index(description, ". "); dotIdx != -1 {
		suffix = description[dotIdx+1:]
	}

	if suffix != "" {
		return prefix + "The ID of the " + resourceType + "." + suffix
	}
	return prefix + "The ID of the " + resourceType + "."
}

// applyFix reads current line from argsSection and applies fix
func (s S012) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, _ string) {
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
		currentLine := content[lineIdx]
		fixedLine := s.fixIdDescription(docProperty.Name, currentLine)
		if fixedLine != "" {
			content[lineIdx] = fixedLine
			argsSection.SetContent(content)
			d.Document.HasChange = true
		}
	}
}
