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

func TestS012_Run(t *testing.T) {
	testCases := []struct {
		name           string
		schemaProps    *models.SchemaProperties
		docProps       *models.DocumentProperties
		expectedErrors int
		errorContains  []string
	}{
		{
			name: "valid id description",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["subnet_id"] = &models.SchemaProperty{Name: "subnet_id", Required: true}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["subnet_id"] = &models.DocumentProperty{
					Name:     "subnet_id",
					Content:  "* `subnet_id` - (Required) The ID of the subnet.",
					Required: true,
				}
				return props
			}(),
			expectedErrors: 0,
		},
		{
			name: "invalid id description",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["subnet_id"] = &models.SchemaProperty{Name: "subnet_id", Required: true}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["subnet_id"] = &models.DocumentProperty{
					Name:     "subnet_id",
					Content:  "* `subnet_id` - (Required) Specifies the subnet.",
					Required: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"should start with 'The ID of the"},
		},
		{
			name: "non-id field skipped",
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["name"] = &models.SchemaProperty{Name: "name", Required: true}
				return props
			}(),
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:     "name",
					Content:  "* `name` - (Required) The name.",
					Required: true,
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
				Type:              data.ResourceTypeResource,
				SchemaProperties:  tc.schemaProps,
				DocumentArguments: tc.docProps,
				Document:          &markdown.Document{Path: "test.markdown", Exists: true},
			}

			rule := S012{}
			errs := rule.Run(nodeData, false)

			if len(errs) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tc.expectedErrors, len(errs), errs)
				return
			}

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
		})
	}
}

func TestS012_FixIdDescription(t *testing.T) {
	testCases := []struct {
		name      string
		fieldName string
		input     string
		expected  string
	}{
		{
			name:      "simple fix",
			fieldName: "subnet_id",
			input:     "* `subnet_id` - (Required) Specifies the subnet.",
			expected:  "* `subnet_id` - (Required) The ID of the subnet.",
		},
		{
			name:      "preserve suffix",
			fieldName: "topic_id",
			input:     "* `topic_id` - (Required) Specifies the topic. Changing this forces a new resource to be created.",
			expected:  "* `topic_id` - (Required) The ID of the topic. Changing this forces a new resource to be created.",
		},
		{
			name:      "optional field",
			fieldName: "virtual_network_id",
			input:     "* `virtual_network_id` - (Optional) The virtual network.",
			expected:  "* `virtual_network_id` - (Optional) The ID of the virtual network.",
		},
	}

	rule := S012{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.fixIdDescription(tc.fieldName, tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}
