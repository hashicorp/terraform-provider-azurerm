// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDiffSuppressKeyVaultVersionedKey(t *testing.T) {
	versionlessKey := "https://example.vault.azure.net/keys/test-key"
	versionedKey := versionlessKey + "/00000000000000000000000000000001"
	otherVersionedKey := "https://example.vault.azure.net/keys/other-key/00000000000000000000000000000001"

	testCases := []struct {
		name                string
		oldValue            string
		newValue            string
		autoRotationEnabled bool
		expected            bool
	}{
		{
			name:                "versionless key with auto rotation",
			oldValue:            versionedKey,
			newValue:            versionlessKey,
			autoRotationEnabled: true,
			expected:            true,
		},
		{
			name:                "versionless key without auto rotation",
			oldValue:            versionedKey,
			newValue:            versionlessKey,
			autoRotationEnabled: false,
			expected:            false,
		},
		{
			name:                "versioned key with auto rotation",
			oldValue:            versionedKey,
			newValue:            versionedKey + "2",
			autoRotationEnabled: true,
			expected:            false,
		},
		{
			name:                "different versionless key with auto rotation",
			oldValue:            otherVersionedKey,
			newValue:            versionlessKey,
			autoRotationEnabled: true,
			expected:            false,
		},
		{
			name:                "invalid configured key",
			oldValue:            versionedKey,
			newValue:            "invalid",
			autoRotationEnabled: true,
			expected:            false,
		},
		{
			name:                "invalid state key",
			oldValue:            "invalid",
			newValue:            versionlessKey,
			autoRotationEnabled: true,
			expected:            false,
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"auto_rotation_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"auto_rotation_enabled": testCase.autoRotationEnabled,
			})

			actual := diffSuppressKeyVaultVersionedKey("", testCase.oldValue, testCase.newValue, data)
			if actual != testCase.expected {
				t.Fatalf("expected diff suppression to be %t, got %t", testCase.expected, actual)
			}
		})
	}

	if diffSuppressKeyVaultVersionedKey("", versionedKey, versionlessKey, nil) {
		t.Fatal("expected diff suppression to be false without resource data")
	}
}
