// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
)

func TestDiskPoolManagedDiskAttachmentIDFormatter(t *testing.T) {
	diskPoolId := diskpools.NewDiskPoolID("12345678-1234-9876-4563-123456789012", "resGroup1", "storagePool1")
	managedDiskId := disks.NewDiskID("12345678-1234-9876-4563-123456789012", "resGroup1", "diks1")
	actual := NewDiskPoolManagedDiskAttachmentId(diskPoolId, managedDiskId).ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/storagePool1/managedDisks|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/diks1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}
