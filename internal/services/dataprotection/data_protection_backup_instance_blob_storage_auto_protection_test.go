// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-03-01/backupinstanceresources"
)

func TestExpandBlobBackupAutoProtection(t *testing.T) {
	withoutRules := expandBlobBackupAutoProtection([]interface{}{})
	settings := withoutRules.AutoProtectionSettings
	if !settings.Enabled {
		t.Fatalf("expected auto protection to be enabled")
	}
	if settings.Rules != nil {
		t.Fatalf("expected no rules when no prefixes are excluded, got %+v", *settings.Rules)
	}

	withRules := expandBlobBackupAutoProtection([]interface{}{"temp-", "test-"})
	settings = withRules.AutoProtectionSettings
	expected := []backupinstanceresources.BlobBackupAutoProtectionRule{
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "temp-"},
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "test-"},
	}
	if !reflect.DeepEqual(pointer.From(settings.Rules), expected) {
		t.Fatalf("unexpected rules:\n want: %+v\n got:  %+v", expected, pointer.From(settings.Rules))
	}
}

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
