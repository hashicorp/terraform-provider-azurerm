// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

func TestS007_Run(t *testing.T) {
	testCases := []struct {
		name           string
		schemaProps    *models.SchemaProperties
		docProps       *models.DocumentProperties
		expectedErrors int
		errorContains  []string
	}{
		{
			name: "requiredness matches",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{
					Name:     "name",
					Required: true,
				}
				props.Objects["tags"] = &models.SchemaProperty{
					Name:     "tags",
					Optional: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:     "name",
					Content:  "* `name` - (Required) The name.",
					Required: true,
				}
				props.Objects["tags"] = &models.DocumentProperty{
					Name:     "tags",
					Content:  "* `tags` - (Optional) Tags.",
					Optional: true,
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name: "should be Required but marked Optional",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["location"] = &models.SchemaProperty{
					Name:     "location",
					Required: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["location"] = &models.DocumentProperty{
					Name:     "location",
					Content:  "* `location` - (Optional) The location.",
					Optional: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should be marked as Required", "location"},
		},
		{
			name: "should be Optional but marked Required",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["description"] = &models.SchemaProperty{
					Name:     "description",
					Optional: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["description"] = &models.DocumentProperty{
					Name:     "description",
					Content:  "* `description` - (Required) The description.",
					Required: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should be marked as Optional", "description"},
		},
		{
			name: "nested property requiredness mismatch",
			schemaProps: func() *models.SchemaProperties {
				nested := models.NewSchemaProperties()
				nested.Objects["nested_field"] = &models.SchemaProperty{
					Name:     "nested_field",
					Required: true,
				}
				props := models.NewSchemaProperties()
				props.Objects["config"] = &models.SchemaProperty{
					Name:     "config",
					Optional: true,
					Nested:   nested,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				nestedDoc := models.NewDocumentProperties()
				nestedDoc.Objects["nested_field"] = &models.DocumentProperty{
					Name:     "nested_field",
					Content:  "* `nested_field` - (Optional) Nested field.",
					Optional: true,
				}
				props := models.NewDocumentProperties()
				props.Objects["config"] = &models.DocumentProperty{
					Name:          "config",
					Optional:      true,
					Block:         true,
					BlockTypeName: "config",
					Nested:        nestedDoc,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should be marked as Required", "config.nested_field"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeData := &data.TerraformNodeData{
				Name:              "azurerm_test_resource",
				Type:              data.ResourceTypeResource,
				SchemaProperties:  tc.schemaProps,
				DocumentArguments: tc.docProps,
				Document: &markdown.Document{
					Path:   "test.markdown",
					Exists: true,
				},
			}

			rule := S007{}
			errs := rule.Run(nodeData, false)

			if len(errs) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tc.expectedErrors, len(errs), errs)
				return
			}

			if tc.expectedErrors > 0 {
				for _, errMsg := range tc.errorContains {
					found := false
					for _, err := range errs {
						if strings.Contains(err.Error(), errMsg) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error to contain '%s', but not found in: %v", errMsg, errs)
					}
				}
			}
		})
	}
}

func TestS007_ReplaceRequiredness(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		from     string
		to       string
		expected string
	}{
		{
			name:     "replace Optional with Required",
			input:    "* `name` - (Optional) The name.",
			from:     "(Optional)",
			to:       "(Required)",
			expected: "* `name` - (Required) The name.",
		},
		{
			name:     "replace Required with Optional",
			input:    "* `location` - (Required) The location.",
			from:     "(Required)",
			to:       "(Optional)",
			expected: "* `location` - (Optional) The location.",
		},
		{
			name:     "add Required when missing",
			input:    "* `name` - The name.",
			from:     "(Optional)",
			to:       "(Required)",
			expected: "* `name` - (Required) The name.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := S007{}
			result := rule.replaceRequiredness(tc.input, tc.from, tc.to)
			if result != tc.expected {
				t.Errorf("Expected:\n  %s\nGot:\n  %s", tc.expected, result)
			}
		})
	}
}
