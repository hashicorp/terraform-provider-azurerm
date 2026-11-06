// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"reflect"
	"testing"
)

func TestDataProtectionBackupVaultV0ToV1(t *testing.T) {
	tests := []struct {
		datastoreType string
		redundancy    string
		expected      []any
	}{
		{
			datastoreType: "VaultStore",
			redundancy:    "GeoRedundant",
			expected: []any{
				map[string]any{"datastore_type": "VaultStore", "redundancy": "GeoRedundant"},
			},
		},
		{
			datastoreType: "OperationalStore",
			redundancy:    "LocallyRedundant",
			expected: []any{
				map[string]any{"datastore_type": "OperationalStore", "redundancy": "LocallyRedundant"},
			},
		},
		{
			datastoreType: "ArchiveStore",
			redundancy:    "LocallyRedundant",
			expected: []any{
				map[string]any{"datastore_type": "ArchiveStore", "redundancy": "LocallyRedundant"},
				map[string]any{"datastore_type": "VaultStore", "redundancy": "LocallyRedundant"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.datastoreType, func(t *testing.T) {
			input := map[string]any{
				"id":                         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.DataProtection/backupVaults/test",
				"datastore_type":             test.datastoreType,
				"redundancy":                 test.redundancy,
				"retention_duration_in_days": float64(30),
				"tags":                       map[string]any{"environment": "test"},
			}
			expected := map[string]any{
				"id":                         input["id"],
				"storage_setting":            test.expected,
				"retention_duration_in_days": float64(30),
				"tags":                       map[string]any{"environment": "test"},
			}

			actual, err := (DataProtectionBackupVaultV0ToV1{}).UpgradeFunc()(context.Background(), input, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("expected %#v, got %#v", expected, actual)
			}
		})
	}
}

func TestDataProtectionBackupVaultV0ToV1_missingFields(t *testing.T) {
	for _, field := range []string{"datastore_type", "redundancy"} {
		t.Run(field, func(t *testing.T) {
			input := map[string]any{"datastore_type": "VaultStore", "redundancy": "LocallyRedundant"}
			input[field] = nil
			if _, err := (DataProtectionBackupVaultV0ToV1{}).UpgradeFunc()(context.Background(), input, nil); err == nil {
				t.Fatal("expected an error for a missing storage setting")
			}
		})
	}
}
