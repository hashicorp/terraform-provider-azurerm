// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/mdparser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

func TestS009_Run(t *testing.T) {
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
			name:           "property missing in documentation",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{
					Name:     "name",
					Required: true,
					Optional: false,
				}
				return props
			}(),
			docProps:       models.NewDocumentProperties(),
			expectedErrors: 1,
			errorContains:  []string{"exists in schema but is MISSING from documentation", "name"},
		},
		{
			name:           "property missing in schema",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps:    models.NewSchemaProperties(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["location"] = &models.DocumentProperty{
					Name:     "location",
					Content:  "* `location` - (Required) The location.",
					Required: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"exists in documentation but NOT in schema", "location"},
		},
		{
			name:           "deprecated property skipped",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["old_field"] = &models.SchemaProperty{
					Name:       "old_field",
					Optional:   true,
					Deprecated: true,
				}
				return props
			}(),
			docProps:       models.NewDocumentProperties(),
			expectedErrors: 0,
		},
		{
			name:           "block not marked in documentation",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				nested := models.NewSchemaProperties()
				nested.Objects["sub_field"] = &models.SchemaProperty{
					Name:     "sub_field",
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
				props := models.NewDocumentProperties()
				props.Objects["config"] = &models.DocumentProperty{
					Name:     "config",
					Content:  "* `config` - (Optional) Configuration.",
					Optional: true,
					Block:    false,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should be declared as a BLOCK"},
		},
		{
			name:           "block incorrectly marked in documentation",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{
					Name:     "name",
					Required: true,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:          "name",
					Content:       "* `name` - (Required) Name.",
					Required:      true,
					Block:         true,
					BlockTypeName: "name",
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"incorrectly implies", "is a BLOCK"},
		},
		{
			name:           "nested properties match",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				nested := models.NewSchemaProperties()
				nested.Objects["nested_field"] = &models.SchemaProperty{
					Name:     "nested_field",
					Required: true,
				}
				props := models.NewSchemaProperties()
				props.Objects["block"] = &models.SchemaProperty{
					Name:     "block",
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
				}
				props := models.NewDocumentProperties()
				props.Objects["block"] = &models.DocumentProperty{
					Name:          "block",
					Content:       "* `block` - (Optional) Block.",
					Optional:      true,
					Block:         true,
					BlockTypeName: "block",
					Nested:        nestedDoc,
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name:           "nested property missing",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				nested := models.NewSchemaProperties()
				nested.Objects["field1"] = &models.SchemaProperty{Name: "field1", Required: true}
				nested.Objects["field2"] = &models.SchemaProperty{Name: "field2", Required: true}
				props := models.NewSchemaProperties()
				props.Objects["block"] = &models.SchemaProperty{
					Name:     "block",
					Optional: true,
					Nested:   nested,
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				nestedDoc := models.NewDocumentProperties()
				nestedDoc.Objects["field1"] = &models.DocumentProperty{
					Name:     "field1",
					Required: true,
				}
				props := models.NewDocumentProperties()
				props.Objects["block"] = &models.DocumentProperty{
					Name:          "block",
					Optional:      true,
					Block:         true,
					BlockTypeName: "block",
					Nested:        nestedDoc,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"block.field2"},
		},
		{
			name:           "misspelling detection",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps:    models.NewSchemaProperties(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				prop := &models.DocumentProperty{
					Name:        "locaton",
					Content:     "* `locaton` - (Required) The location.",
					Required:    true,
					ParseErrors: []string{mdparser.MisspelNameOfProperty + " from `location` to `locaton`"},
				}
				props.Objects["locaton"] = prop
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"misspelling", "locaton"},
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

			rule := S009{}
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
