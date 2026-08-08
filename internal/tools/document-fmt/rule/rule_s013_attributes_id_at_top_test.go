// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

func TestS013_Run(t *testing.T) {
	testCases := []struct {
		name           string
		content        []string
		expectedErrors int
		errorContains  []string
	}{
		{
			name: "id is first attribute",
			content: []string{
				"## Attributes Reference",
				"",
				"* `id` - The ID of the resource.",
				"* `name` - The name.",
			},
			expectedErrors: 0,
		},
		{
			name: "id is not first attribute",
			content: []string{
				"## Attributes Reference",
				"",
				"* `name` - The name.",
				"* `id` - The ID of the resource.",
			},
			expectedErrors: 1,
			errorContains:  []string{"should be the first attribute"},
		},
		{
			name: "no id attribute",
			content: []string{
				"## Attributes Reference",
				"",
				"* `name` - The name.",
			},
			expectedErrors: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attrSection := &markdown.AttributesSection{}
			attrSection.SetContent(tc.content)

			nodeData := &data.TerraformNodeData{
				Name: "azurerm_test_resource",
				Type: data.ResourceTypeResource,
				Document: &markdown.Document{
					Path:     "test.markdown",
					Exists:   true,
					Sections: []markdown.Section{attrSection},
				},
			}

			rule := S013{}
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
