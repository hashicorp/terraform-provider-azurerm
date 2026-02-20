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

func TestS006_Run(t *testing.T) {
	testCases := []struct {
		name           string
		docProps       *models.DocumentProperties
		schemaProps    *models.SchemaProperties
		expectedErrors int
		errorContains  []string
	}{
		{
			name: "no parse errors",
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:     "name",
					Content:  "* `name` - (Required) The name.",
					Required: true,
				}
				return props
			}(),
			schemaProps:    models.NewSchemaProperties(),
			expectedErrors: 0,
		},
		{
			name: "incorrectly marked as block",
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["config"] = &models.DocumentProperty{
					Name:        "config",
					Content:     "* `config` - (Optional) Configuration as defined below.",
					Optional:    true,
					ParseErrors: []string{mdparser.IncorrectlyBlockMarked},
				}
				return props
			}(),
			schemaProps: func() *models.SchemaProperties {
				props := models.NewSchemaProperties()
				props.Objects["config"] = &models.SchemaProperty{
					Name:     "config",
					Optional: true,
				}
				return props
			}(),
			expectedErrors: 1,
			errorContains:  []string{"incorrectly implies", "is a block"},
		},
		{
			name: "duplicate property",
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["name"] = &models.DocumentProperty{
					Name:        "name",
					Content:     "* `name` - (Required) The name.",
					Required:    true,
					ParseErrors: []string{mdparser.DuplicateFieldsFound},
				}
				return props
			}(),
			schemaProps:    models.NewSchemaProperties(),
			expectedErrors: 1,
			errorContains:  []string{"Duplicate fields declared"},
		},
		{
			name: "no field name found",
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["bad"] = &models.DocumentProperty{
					Name:        "bad",
					Required:    true,
					Content:     "* Some malformed line without proper format",
					ParseErrors: []string{mdparser.NoFieldNameFound},
				}
				return props
			}(),
			schemaProps:    models.NewSchemaProperties(),
			expectedErrors: 1,
			errorContains:  []string{"should be formatted as"},
		},
		{
			name: "nested property with format error",
			docProps: func() *models.DocumentProperties {
				nestedDoc := models.NewDocumentProperties()
				nestedDoc.Objects["nested_field"] = &models.DocumentProperty{
					Name:        "nested_field",
					Required:    true,
					Content:     "* `nested_field` - Field as defined below.",
					ParseErrors: []string{mdparser.IncorrectlyBlockMarked},
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
			expectedErrors: 1,
			errorContains:  []string{"incorrectly implies", "is a block"},
		},
		{
			name: "misspelling error skipped",
			docProps: func() *models.DocumentProperties {
				props := models.NewDocumentProperties()
				props.Objects["locaton"] = &models.DocumentProperty{
					Name:        "locaton",
					Optional:    true,
					Content:     "* `locaton` - (Required) Location.",
					ParseErrors: []string{mdparser.MisspelNameOfProperty},
				}
				return props
			}(),
			schemaProps:    models.NewSchemaProperties(),
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
				Document: &markdown.Document{
					Path:   "test.markdown",
					Exists: true,
				},
			}

			rule := S006{}
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

func TestS006_FixFormatError(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		parseErr string
		expected string
	}{
		{
			name:     "remove star from note line",
			input:    "* ~> **NOTE:** This is a note.",
			parseErr: "some error",
			expected: "~> **NOTE:** This is a note.",
		},
		{
			name:     "clear empty list marker",
			input:    "*",
			parseErr: "empty marker",
			expected: "",
		},
		{
			name:     "line with incorrectly block marked passes through",
			input:    "* `config` - (Optional) Configuration as defined below.",
			parseErr: "incorrectly block marked",
			expected: "* `config` - (Optional) Configuration as defined below.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := S006{}
			result := rule.fixFormatError(tc.input, tc.parseErr)
			if result != tc.expected {
				t.Errorf("Expected:\n  %s\nGot:\n  %s", tc.expected, result)
			}
		})
	}
}
