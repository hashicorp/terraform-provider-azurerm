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

func TestS010_Run(t *testing.T) {
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
			name:           "default values match",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["enabled"] = &models.SchemaProperty{
					Name:         "enabled",
					Optional:     true,
					Type:         "Bool",
					DefaultValue: "true",
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["enabled"] = &models.DocumentProperty{
					Name:         "enabled",
					Content:      "* `enabled` - (Optional) Whether enabled. Defaults to `true`.",
					Optional:     true,
					DefaultValue: "true",
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name:           "schema has default, doc missing",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["timeout"] = &models.SchemaProperty{
					Name:         "timeout",
					Optional:     true,
					Type:         "Int",
					DefaultValue: "300",
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["timeout"] = &models.DocumentProperty{
					Name:     "timeout",
					Content:  "* `timeout` - (Optional) The timeout in seconds.",
					Optional: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"missing / wrong Default Value declaration", "timeout"},
		},
		{
			name:           "default values mismatch",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["retry_count"] = &models.SchemaProperty{
					Name:         "retry_count",
					Optional:     true,
					Type:         "Int",
					DefaultValue: "3",
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["retry_count"] = &models.DocumentProperty{
					Name:         "retry_count",
					Content:      "* `retry_count` - (Optional) Retry count. Defaults to `5`.",
					Optional:     true,
					DefaultValue: "5",
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"missing / wrong Default Value declaration", "retry_count"},
		},
		{
			name:           "doc has default but schema doesn't (non-bool)",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{
					Name:     "name",
					Optional: true,
					Type:     "String",
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:         "name",
					Content:      "* `name` - (Optional) The name. Defaults to `default`.",
					Optional:     true,
					DefaultValue: "default",
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should NOT have default value", "name"},
		},
		{
			name:           "false default for bool - doc omitted",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["enabled"] = &models.SchemaProperty{
					Name:         "enabled",
					Optional:     true,
					Type:         "Bool",
					DefaultValue: "false",
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["enabled"] = &models.DocumentProperty{
					Name:     "enabled",
					Content:  "* `enabled` - (Optional) Whether enabled.",
					Optional: true,
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

			rule := S010{}
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
