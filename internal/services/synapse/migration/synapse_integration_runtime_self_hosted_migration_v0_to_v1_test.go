// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestSynapseIntegrationRuntimeSelfHostedV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "old id - lowercased static segments",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/group1/providers/microsoft.synapse/workspaces/workspace1/integrationruntimes/runtime1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/integrationRuntimes/runtime1"),
		},
		{
			name: "old id - uppercased static segments",
			input: map[string]interface{}{
				"id": "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/group1/PROVIDERS/MICROSOFT.SYNAPSE/WORKSPACES/workspace1/INTEGRATIONRUNTIMES/runtime1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/integrationRuntimes/runtime1"),
		},
		{
			name: "old id - user specified segment casing is preserved",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/Group1/providers/microsoft.synapse/workspaces/Workspace1/integrationruntimes/Runtime1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/Group1/providers/Microsoft.Synapse/workspaces/Workspace1/integrationRuntimes/Runtime1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/integrationRuntimes/runtime1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1/integrationRuntimes/runtime1"),
		},
		{
			name: "invalid id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Synapse/workspaces/workspace1",
			},
			expected: nil,
		},
	}

	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		result, err := SynapseIntegrationRuntimeSelfHostedV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
		if err != nil {
			if test.expected == nil {
				continue
			}
			t.Fatalf("Expected no error but got: %+v", err)
		}
		if test.expected == nil {
			t.Fatalf("Expected an error but didn't get one")
		}

		actualId := result["id"].(string)
		if *test.expected != actualId {
			t.Fatalf("expected %q but got %q!", *test.expected, actualId)
		}
	}
}
