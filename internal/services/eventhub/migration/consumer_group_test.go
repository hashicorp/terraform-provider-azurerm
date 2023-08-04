// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestConsumerGroupsV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "old id (shared)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumergroup1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumergroup1"),
		},
		{
			name: "old id - mixed case (shared)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/ConsumerGroups/consumergroup1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumergroup1"),
		},
		{
			name: "new id (shared)",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumergroup1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumerGroups/consumergroup1"),
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		result, err := ConsumerGroupsV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
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
