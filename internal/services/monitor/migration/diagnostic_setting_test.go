package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestDiagnosticSettingV0ToV1(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "old id",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"),
		},
		{
			name: "old id - name",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"),
		},
		{
			name: "old id - mix",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"),
		},
		{
			name: "old id - mix",
			input: map[string]interface{}{
				"id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourcegroups/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1"),
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q...", test.name)
		result, err := DiagnosticSettingUpgradeV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
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
