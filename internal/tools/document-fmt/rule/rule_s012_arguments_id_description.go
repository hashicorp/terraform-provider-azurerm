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

	// Extract description after "(Required)" or "(Optional)"
	content := docProperty.Content
	idx := strings.Index(content, ") ")
	if idx == -1 {
		return nil
	}
	description := strings.TrimSpace(content[idx+2:])

	if strings.HasPrefix(strings.ToLower(description), "the id of the") {
		return nil
	}

	return NewValidationIssue(s.ID(), s.Name(), fullPath,
		fmt.Sprintf("`%s` description should start with 'The ID of the ...'", util.Bold(fullPath)),
		d.Document.Path, content, "")
}

