package diskpools

import (
	"testing"

	computeparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func TestDiskPoolManagedDiskAttachmentIDFormatter(t *testing.T) {
	diskPoolId := NewDiskPoolID("12345678-1234-9876-4563-123456789012", "resGroup1", "storagePool1")
	managedDiskId := computeparse.NewManagedDiskID("12345678-1234-9876-4563-123456789012", "resGroup1", "diks1")
	actual := NewDiskPoolManagedDiskAttachmentId(diskPoolId, managedDiskId).ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/storagePool1/managedDisks|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/diks1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}
