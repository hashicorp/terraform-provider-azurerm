// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
)

func TestFlattenBlobBackupAutoProtection(t *testing.T) {
	testCases := []struct {
		name             string
		input            backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection
		expectedEnabled  bool
		expectedPrefixes []string
	}{
		{
			name:             "zero value",
			input:            backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection{},
			expectedEnabled:  false,
			expectedPrefixes: []string{},
		},
		{
			name: "enabled without rules",
			input: backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection{
				AutoProtectionSettings: backupinstanceresources.BlobBackupRuleBasedAutoProtectionSettings{Enabled: true},
			},
			expectedEnabled:  true,
			expectedPrefixes: []string{},
		},
		{
			name: "enabled with rules keeps order and ignores unknown rule kinds",
			input: backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection{
				AutoProtectionSettings: backupinstanceresources.BlobBackupRuleBasedAutoProtectionSettings{
					Enabled: true,
					Rules: &[]backupinstanceresources.BlobBackupAutoProtectionRule{
						{Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "temp-"},
						{Mode: "Include", Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "keep-"},
						{Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "test-"},
					},
				},
			},
			expectedEnabled:  true,
			expectedPrefixes: []string{"temp-", "test-"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enabled, prefixes := flattenBlobBackupAutoProtection(tc.input)
			if enabled != tc.expectedEnabled {
				t.Fatalf("expected enabled = %t, got %t", tc.expectedEnabled, enabled)
			}
			if !reflect.DeepEqual(prefixes, tc.expectedPrefixes) {
				t.Fatalf("expected prefixes %v, got %v", tc.expectedPrefixes, prefixes)
			}
		})
	}
}
