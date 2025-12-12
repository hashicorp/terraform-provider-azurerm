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

func TestS011_Run(t *testing.T) {
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
			name:           "possible values match",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["frequency"] = &models.SchemaProperty{
					Name:           "frequency",
					Required:       true,
					Type:           "String",
					PossibleValues: []string{"Daily", "Weekly"},
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["frequency"] = &models.DocumentProperty{
					Name:      "frequency",
					Content:   "* `frequency` - (Required) The backup frequency. Possible values are `Daily` and `Weekly`.",
					Required:  true,
					EnumStart: 59,
					EnumEnd:   82,
					Enums:     map[string]struct{}{"Daily": {}, "Weekly": {}},
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name:           "schema has possible values, doc missing",
			resourceType:   data.ResourceTypeResource,
			documentExists: true,
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["sku"] = &models.SchemaProperty{
					Name:           "sku",
					Required:       true,
					Type:           "String",
					PossibleValues: []string{"Basic", "Standard", "Premium"},
				}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["sku"] = &models.DocumentProperty{
					Name:     "sku",
					Content:  "* `sku` - (Required) The SKU name.",
					Required: true,
					Enums:    map[string]struct{}{},
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"missing possible values", "sku"},
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

			rule := S011{}
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
