// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
)

func TestExpandBlobBackupAutoProtection(t *testing.T) {
	withoutRules := expandBlobBackupAutoProtection([]any{})
	settings := withoutRules.AutoProtectionSettings
	if !settings.Enabled {
		t.Fatalf("expected auto protection to be enabled")
	}
	if settings.Rules != nil {
		t.Fatalf("expected no rules when no prefixes are excluded, got %+v", *settings.Rules)
	}

	withRules := expandBlobBackupAutoProtection([]any{"temp-", "test-"})
	settings = withRules.AutoProtectionSettings
	expected := []backupinstanceresources.BlobBackupAutoProtectionRule{
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "temp-"},
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: backupinstanceresources.BlobBackupRuleModeExclude, Type: backupinstanceresources.BlobBackupPatternTypePrefix, Pattern: "test-"},
	}
	if !reflect.DeepEqual(pointer.From(settings.Rules), expected) {
		t.Fatalf("unexpected rules:\n want: %+v\n got:  %+v", expected, pointer.From(settings.Rules))
	}
}
