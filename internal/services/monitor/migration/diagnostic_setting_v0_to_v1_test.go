// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestDiagnosticSettingV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]any
		expected *string
	}{
		{
			name: "legacy pipe id",
			input: map[string]any{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1|setting1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1/providers/Microsoft.Insights/diagnosticSettings/setting1"),
		},
		{
			name: "already arm id",
			input: map[string]any{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1/providers/Microsoft.Insights/diagnosticSettings/setting1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1/providers/Microsoft.Insights/diagnosticSettings/setting1"),
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			result, err := DiagnosticSettingV0ToV1{}.UpgradeFunc()(context.Background(), test.input, nil)
			if err != nil && test.expected != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if test.expected != nil {
				actual := result["id"].(string)
				if actual != *test.expected {
					t.Fatalf("expected ID %q, got %q", *test.expected, actual)
				}
			}
		})
	}
}
