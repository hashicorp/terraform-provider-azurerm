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

func TestS008_Run(t *testing.T) {
	testCases := []struct {
		name           string
		resourceType   data.ResourceType
		documentExists bool
		schemaProps    *models.SchemaProperties
		docProps       *models.DocumentProperties
		expectedErrors int
		errorContains  []string
	}{
		{
			name:           "ForceNew matches",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{
					Name:     "name",
					Required: true,
					ForceNew: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:     "name",
					Content:  "* `name` - (Required) The name. Changing this forces a new resource to be created.",
					Required: true,
					ForceNew: true,
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name:           "ForceNew missing in doc",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["location"] = &models.SchemaProperty{
					Name:     "location",
					Required: true,
					ForceNew: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["location"] = &models.DocumentProperty{
					Name:     "location",
					Content:  "* `location` - (Required) The location.",
					Required: true,
					ForceNew: false,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should be marked as ForceNew", "location"},
		},
		{
			name:           "ForceNew should not be in doc",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["tags"] = &models.SchemaProperty{
					Name:     "tags",
					Optional: true,
					ForceNew: false,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["tags"] = &models.DocumentProperty{
					Name:     "tags",
					Content:  "* `tags` - (Optional) Tags. Changing this forces a new resource to be created.",
					Optional: true,
					ForceNew: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should NOT be marked as ForceNew", "tags"},
		},
		{
			name:           "nested property ForceNew mismatch",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				nested := models.NewSchemaProperties()
				nested.Objects["nested_field"] = &models.SchemaProperty{
					Name:     "nested_field",
					Required: true,
					ForceNew: true,
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
					Content:  "* `nested_field` - (Required) Nested field.",
					Required: true,
					ForceNew: false,
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
			errorContains:  []string{"should be marked as ForceNew", "config.nested_field"},
		},
		{
			name:           "resource_group_name skipped",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["resource_group_name"] = &models.SchemaProperty{
					Name:     "resource_group_name",
					Required: true,
					ForceNew: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["resource_group_name"] = &models.DocumentProperty{
					Name:     "resource_group_name",
					Content:  "* `resource_group_name` - (Required) The name of the resource group.",
					Required: true,
					ForceNew: false,
				}
				return props
			}(),
			expectedErrors: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeData := &data.TerraformNodeData{
				Name:              "azurerm_test_resource",
				Type:              tc.resourceType,
				SchemaProperties:  tc.schemaProps,
				DocumentArguments: tc.docProps,
				Document: &markdown.Document{
					Path:   "test.markdown",
					Exists: tc.documentExists,
				},
			}

			rule := S008{}
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

func TestS008_FixForceNew(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		add      bool
		expected string
	}{
		{
			name:     "add ForceNew to line ending with period",
			input:    "* `name` - (Required) The name.",
			add:      true,
			expected: "* `name` - (Required) The name. Changing this forces a new resource to be created.",
		},
		{
			name:     "add ForceNew to line without period",
			input:    "* `name` - (Required) The name",
			add:      true,
			expected: "* `name` - (Required) The name. Changing this forces a new resource to be created.",
		},
		{
			name:     "add ForceNew to line ending with comma",
			input:    "* `name` - (Required) The name,",
			add:      true,
			expected: "* `name` - (Required) The name. Changing this forces a new resource to be created.",
		},
		{
			name:     "remove ForceNew from line",
			input:    "* `name` - (Required) The name. Changing this forces a new resource to be created.",
			add:      false,
			expected: "* `name` - (Required) The name.",
		},
		{
			name:     "remove ForceNew variation",
			input:    "* `location` - (Required) Location. Changing this will force a new resource to be created.",
			add:      false,
			expected: "* `location` - (Required) Location.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := S008{}
			result := rule.fixForceNew(tc.input, tc.add)
			if result != tc.expected {
				t.Errorf("Expected:\n  %s\nGot:\n  %s", tc.expected, result)
			}
		})
	}
}

func TestS008_LastPathSegment(t *testing.T) {
	testCases := []struct {
		path     string
		expected string
	}{
		{"name", "name"},
		{"block.name", "name"},
		{"outer.inner.field", "field"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := lastPathSegment(tc.path)
			if result != tc.expected {
				t.Errorf("lastPathSegment(%s) = %s, expected %s", tc.path, result, tc.expected)
			}
		})
	}
}
