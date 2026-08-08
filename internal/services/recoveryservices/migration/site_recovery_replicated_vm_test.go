// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestSiteRecoveryReplicatedVMV0ToV1(t *testing.T) {
	// `managed_disk` switched from a TypeSet to a TypeList. Both are persisted as a JSON array, so the
	// upgrade is a no-op on the value itself - the version bump is what re-interprets it as a list. The
	// upgrade must preserve the disks (order and contents) so Read can re-align them to config order.
	input := map[string]interface{}{
		"name": "repl-1",
		"managed_disk": []interface{}{
			map[string]interface{}{
				"disk_id":                    "/subscriptions/s/resourceGroups/rg/providers/Microsoft.Compute/disks/disk1",
				"staging_storage_account_id": "/subscriptions/s/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sa1",
				"target_resource_group_id":   "/subscriptions/s/resourceGroups/rg",
				"target_disk_type":           "Premium_LRS",
				"target_replica_disk_type":   "Premium_LRS",
			},
			map[string]interface{}{
				"disk_id":                    "/subscriptions/s/resourceGroups/rg/providers/Microsoft.Compute/disks/disk2",
				"staging_storage_account_id": "/subscriptions/s/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sa1",
				"target_resource_group_id":   "/subscriptions/s/resourceGroups/rg",
				"target_disk_type":           "StandardSSD_LRS",
				"target_replica_disk_type":   "StandardSSD_LRS",
			},
		},
	}

	expected := map[string]interface{}{
		"name":         "repl-1",
		"managed_disk": input["managed_disk"],
	}

	got, err := SiteRecoveryReplicatedVMV0ToV1{}.UpgradeFunc()(context.Background(), input, nil)
	if err != nil {
		t.Fatalf("upgrading state: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected upgraded state:\n got: %#v\nwant: %#v", got, expected)
	}
}

func TestSiteRecoveryReplicatedVMV0SchemaManagedDiskIsSet(t *testing.T) {
	managedDisk, ok := SiteRecoveryReplicatedVMV0ToV1{}.Schema()["managed_disk"]
	if !ok {
		t.Fatal("expected managed_disk in the v0 schema")
	}
	if managedDisk.Type != pluginsdk.TypeSet {
		t.Fatalf("v0 managed_disk must be a TypeSet, got %v", managedDisk.Type)
	}
}
