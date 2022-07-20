package parse

// This is manual for concat two ids are not supported in auto-generation

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2021-04-01/objectreplicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ObjectReplicationId{}

func TestObjectReplicationIDFormatter(t *testing.T) {
	actual := NewObjectReplicationID(
		objectreplicationpolicies.NewObjectReplicationPoliciesID("12345678-1234-9876-4563-123456789012", "resGroup1", "storageAccount1", "objectReplicationPolicy1"),
		objectreplicationpolicies.NewObjectReplicationPoliciesID("12345678-1234-9876-4563-123456789012", "resGroup2", "storageAccount2", "objectReplicationPolicy2"),
	).ID()
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
				Src: objectreplicationpolicies.ObjectReplicationPoliciesId{
					SubscriptionId:            "12345678-1234-9876-4563-123456789012",
					ResourceGroupName:         "resGroup1",
					AccountName:               "storageAccount1",
					ObjectReplicationPolicyId: "objectReplicationPolicy1",
				},
				Dst: objectreplicationpolicies.ObjectReplicationPoliciesId{
					SubscriptionId:            "12345678-1234-9876-4563-123456789012",
					ResourceGroupName:         "resGroup2",
					AccountName:               "storageAccount2",
					ObjectReplicationPolicyId: "objectReplicationPolicy2",
				},
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

		if actual.Src.SubscriptionId != v.Expected.Src.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.Src.SubscriptionId, actual.Src.SubscriptionId)
		}
		if actual.Src.ResourceGroupName != v.Expected.Src.ResourceGroupName {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.Src.ResourceGroupName, actual.Src.ResourceGroupName)
		}
		if actual.Src.AccountName != v.Expected.Src.AccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.Src.AccountName, actual.Src.AccountName)
		}
		if actual.Src.ObjectReplicationPolicyId != v.Expected.Src.ObjectReplicationPolicyId {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Src.ObjectReplicationPolicyId, actual.Src.ObjectReplicationPolicyId)
		}
		if actual.Dst.SubscriptionId != v.Expected.Dst.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.Dst.SubscriptionId, actual.Dst.SubscriptionId)
		}
		if actual.Dst.ResourceGroupName != v.Expected.Dst.ResourceGroupName {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.Dst.ResourceGroupName, actual.Dst.ResourceGroupName)
		}
		if actual.Dst.AccountName != v.Expected.Dst.AccountName {
			t.Fatalf("Expected %q but got %q for StorageAccountName", v.Expected.Dst.AccountName, actual.Dst.AccountName)
		}
		if actual.Dst.ObjectReplicationPolicyId != v.Expected.Dst.ObjectReplicationPolicyId {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Dst.ObjectReplicationPolicyId, actual.Dst.ObjectReplicationPolicyId)
		}
	}
}
