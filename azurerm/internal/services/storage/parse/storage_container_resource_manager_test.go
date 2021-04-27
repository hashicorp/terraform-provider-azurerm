package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = StorageContainerResourceManagerId{}

func TestStorageContainerResourceManagerIDFormatter(t *testing.T) {
	actual := NewStorageContainerResourceManagerID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageAccount1", "default", "container1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default/containers/container1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageContainerResourceManagerID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageContainerResourceManagerId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/",
			Error: true,
		},

		{
			// missing value for StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/",
			Error: true,
		},

		{
			// missing BlobServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Error: true,
		},

		{
			// missing value for BlobServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/",
			Error: true,
		},

		{
			// missing ContainerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default/",
			Error: true,
		},

		{
			// missing value for ContainerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default/containers/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/blobServices/default/containers/container1",
			Expected: &StorageContainerResourceManagerId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				StorageAccountName: "storageAccount1",
				BlobServiceName:    "default",
				ContainerName:      "container1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT1/BLOBSERVICES/DEFAULT/CONTAINERS/CONTAINER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageContainerResourceManagerID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.StorageAccountName != v.Expected.StorageAccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.StorageAccountName, actual.StorageAccountName)
		}
		if actual.BlobServiceName != v.Expected.BlobServiceName {
			t.Fatalf("Expected %q but got %q for BlobServiceName", v.Expected.BlobServiceName, actual.BlobServiceName)
		}
		if actual.ContainerName != v.Expected.ContainerName {
			t.Fatalf("Expected %q but got %q for ContainerName", v.Expected.ContainerName, actual.ContainerName)
		}
	}
}
