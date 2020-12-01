package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = StorageShareResourceManagerId{}

func TestStorageShareResourceManagerIDFormatter(t *testing.T) {
	actual := NewStorageShareResourceManagerID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageAccount1", "fileService1", "share1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/fileService1/shares/share1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageShareResourceManagerID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageShareResourceManagerId
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
			// missing FileServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Error: true,
		},

		{
			// missing value for FileServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/",
			Error: true,
		},

		{
			// missing ShareName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/fileService1/",
			Error: true,
		},

		{
			// missing value for ShareName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/fileService1/shares/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/fileServices/fileService1/shares/share1",
			Expected: &StorageShareResourceManagerId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				StorageAccountName: "storageAccount1",
				FileServiceName:    "fileService1",
				ShareName:          "share1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT1/FILESERVICES/FILESERVICE1/SHARES/SHARE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageShareResourceManagerID(v.Input)
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
		if actual.FileServiceName != v.Expected.FileServiceName {
			t.Fatalf("Expected %q but got %q for FileServiceName", v.Expected.FileServiceName, actual.FileServiceName)
		}
		if actual.ShareName != v.Expected.ShareName {
			t.Fatalf("Expected %q but got %q for ShareName", v.Expected.ShareName, actual.ShareName)
		}
	}
}
