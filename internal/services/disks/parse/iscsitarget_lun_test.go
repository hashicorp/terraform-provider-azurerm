// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
)

func TestNewDiskPoolIscsiTargetLunId(t *testing.T) {
	iscsiTargetId := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")
	managedDiskId := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "resGroup1", "disk1")
	id := NewDiskPoolIscsiTargetLunId(iscsiTargetId, managedDiskId)

	if id.IscsiTargetId != iscsiTargetId {
		t.Fatalf("Expected iscsi taraget id:%s, actual is:%s", iscsiTargetId.ID(), id.IscsiTargetId.ID())
	}

	if id.ManagedDiskId != managedDiskId {
		t.Fatalf("Expected lun is:%s, actual is %s", managedDiskId, id.ManagedDiskId)
	}
}

func TestFormatIscsiTargetLunId(t *testing.T) {
	iscsiTargetId := iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue")
	managedDiskId := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "resGroup1", "disk1")
	id := NewDiskPoolIscsiTargetLunId(iscsiTargetId, managedDiskId)

	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue/lun|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1"

	if expected != id.ID() {
		t.Fatalf("Expected id: %s, actual is:%s", expected, id.ID())
	}
}

func TestParseIscsiTargetLunID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DiskPoolIscsiTargetLunId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue/lun|",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue/lun|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1",
			Expected: &DiskPoolIscsiTargetLunId{
				IscsiTargetId: iscsitargets.NewIscsiTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "diskPoolValue", "iscsiTargetValue"),
				ManagedDiskId: disks.NewDiskID("12345678-1234-9876-4563-123456789012", "resGroup1", "disk1"),
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue/lun|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := IscsiTargetLunID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.IscsiTargetId != v.Expected.IscsiTargetId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.IscsiTargetId, actual.IscsiTargetId)
		}

		if actual.ManagedDiskId != v.Expected.ManagedDiskId {
			t.Fatalf("Expected %q but got %q for ResourceGroupName", v.Expected.ManagedDiskId, actual.ManagedDiskId)
		}
	}
}
