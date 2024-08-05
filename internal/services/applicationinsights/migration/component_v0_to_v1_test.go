// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestComponentV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "old id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/microsoft.insights/components/component1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1"),
		},
		{
			name: "old id - mixed case",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/microsoft.insights/Components/component1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1"),
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		result, err := ComponentUpgradeV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		actualId := result["id"].(string)
		if *test.expected != actualId {
			t.Fatalf("expected %q but got %q!", *test.expected, actualId)
		}
	}
}
