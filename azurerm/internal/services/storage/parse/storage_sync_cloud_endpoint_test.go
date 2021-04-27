package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = StorageSyncCloudEndpointId{}

func TestStorageSyncCloudEndpointIDFormatter(t *testing.T) {
	actual := NewStorageSyncCloudEndpointID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageSyncService1", "syncGroup1", "cloudEndpoint1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/syncGroups/syncGroup1/cloudEndpoints/cloudEndpoint1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStorageSyncCloudEndpointID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageSyncCloudEndpointId
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
			// missing StorageSyncServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/",
			Error: true,
		},

		{
			// missing value for StorageSyncServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/",
			Error: true,
		},

		{
			// missing SyncGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/",
			Error: true,
		},

		{
			// missing value for SyncGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/syncGroups/",
			Error: true,
		},

		{
			// missing CloudEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/syncGroups/syncGroup1/",
			Error: true,
		},

		{
			// missing value for CloudEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/syncGroups/syncGroup1/cloudEndpoints/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/storageSyncService1/syncGroups/syncGroup1/cloudEndpoints/cloudEndpoint1",
			Expected: &StorageSyncCloudEndpointId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "resGroup1",
				StorageSyncServiceName: "storageSyncService1",
				SyncGroupName:          "syncGroup1",
				CloudEndpointName:      "cloudEndpoint1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGESYNC/STORAGESYNCSERVICES/STORAGESYNCSERVICE1/SYNCGROUPS/SYNCGROUP1/CLOUDENDPOINTS/CLOUDENDPOINT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageSyncCloudEndpointID(v.Input)
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
		if actual.StorageSyncServiceName != v.Expected.StorageSyncServiceName {
			t.Fatalf("Expected %q but got %q for StorageSyncServiceName", v.Expected.StorageSyncServiceName, actual.StorageSyncServiceName)
		}
		if actual.SyncGroupName != v.Expected.SyncGroupName {
			t.Fatalf("Expected %q but got %q for SyncGroupName", v.Expected.SyncGroupName, actual.SyncGroupName)
		}
		if actual.CloudEndpointName != v.Expected.CloudEndpointName {
			t.Fatalf("Expected %q but got %q for CloudEndpointName", v.Expected.CloudEndpointName, actual.CloudEndpointName)
		}
	}
}
