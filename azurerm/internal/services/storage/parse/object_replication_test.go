package parse

// This is manual for concat two ids are not supported in auto-generation

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ObjectReplicationId{}

func TestObjectReplicationIDFormatter(t *testing.T) {
	actual := NewObjectReplicationID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageAccount1", "objectReplicationPolicy1", "12345678-1234-9876-4563-123456789012", "resGroup2", "storageAccount2", "objectReplicationPolicy2").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/storageAccount2/objectReplicationPolicies/objectReplicationPolicy2"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestObjectReplicationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ObjectReplicationId
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
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/",
			Error: true,
		},

		{
			// missing SubscriptionName
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionName
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/",
			Error: true,
		},

		{
			// missing value for StorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/storageAccount2/objectReplicationPolicies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/storageAccount2/objectReplicationPolicies/objectReplicationPolicy2",
			Expected: &ObjectReplicationId{
				SrcSubscriptionId:     "12345678-1234-9876-4563-123456789012",
				SrcResourceGroup:      "resGroup1",
				SrcStorageAccountName: "storageAccount1",
				SrcName:               "objectReplicationPolicy1",
				DstSubscriptionId:     "12345678-1234-9876-4563-123456789012",
				DstResourceGroup:      "resGroup2",
				DstStorageAccountName: "storageAccount2",
				DstName:               "objectReplicationPolicy2",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT1/OBJECTREPLICATIONPOLICIES/OBJECTREPLICATIONPOLICY1;/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP2/PROVIDERS/MICROSOFT.STORAGE/STORAGEACCOUNTS/STORAGEACCOUNT2/OBJECTREPLICATIONPOLICIES/OBJECTREPLICATIONPOLICY2",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ObjectReplicationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SrcSubscriptionId != v.Expected.SrcSubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SrcSubscriptionId, actual.SrcSubscriptionId)
		}
		if actual.SrcResourceGroup != v.Expected.SrcResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.SrcResourceGroup, actual.SrcResourceGroup)
		}
		if actual.SrcStorageAccountName != v.Expected.SrcStorageAccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.SrcStorageAccountName, actual.SrcStorageAccountName)
		}
		if actual.SrcName != v.Expected.SrcName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.SrcName, actual.SrcName)
		}
		if actual.DstSubscriptionId != v.Expected.DstSubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionName", v.Expected.DstSubscriptionId, actual.DstSubscriptionId)
		}
		if actual.DstResourceGroup != v.Expected.DstResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.DstResourceGroup, actual.DstResourceGroup)
		}
		if actual.DstStorageAccountName != v.Expected.DstStorageAccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.DstStorageAccountName, actual.DstStorageAccountName)
		}
		if actual.DstName != v.Expected.DstName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.DstName, actual.DstName)
		}
	}
}
