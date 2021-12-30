package parse

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ resourceid.Formatter = StorageAccountNetworkRuleId{}

func TestStorageAccountNetworkRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageAccountNetworkRuleId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing value for StorageAccountID
			Input: "ipAddressOrRange/127.0.0.1",
			Error: true,
		},

		{
			// missing value for networkRuleId
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1",
			Error: true,
		},

		{
			// valid IPRule
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1;ipAddressOrRange/127.0.0.1",
			Expected: &StorageAccountNetworkRuleId{
				StorageAccountId: &StorageAccountId{
					SubscriptionId: "12345678-1234-9876-4563-123456789012",
					ResourceGroup:  "resGroup1",
					Name:           "storageAccount1",
				},
				IPRule: &storage.IPRule{
					IPAddressOrRange: utils.String("127.0.0.1"),
				},
			},
		},

		{
			// valid VirtualNetworkRule
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1;subnetId/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/subnets/subnet1",
			Expected: &StorageAccountNetworkRuleId{
				StorageAccountId: &StorageAccountId{
					SubscriptionId: "12345678-1234-9876-4563-123456789012",
					ResourceGroup:  "resGroup1",
					Name:           "storageAccount1",
				},
				VirtualNetworkRule: &storage.VirtualNetworkRule{
					VirtualNetworkResourceID: utils.String("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/subnets/subnet1"),
				},
			},
		},

		{
			// valid ResourceAccessRule
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storageAccount1;tenantId/12345678-1234-9876-4563-123456789012/resourceId/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/storageAccount2",
			Expected: &StorageAccountNetworkRuleId{
				StorageAccountId: &StorageAccountId{
					SubscriptionId: "12345678-1234-9876-4563-123456789012",
					ResourceGroup:  "resGroup1",
					Name:           "storageAccount1",
				},
				ResourceAccessRule: &storage.ResourceAccessRule{
					TenantID:   utils.String("12345678-1234-9876-4563-123456789012"),
					ResourceID: utils.String("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup2/providers/Microsoft.Storage/storageAccounts/storageAccount2"),
				},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StorageAccountNetworkRuleID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}
			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.StorageAccountId != v.Expected.StorageAccountId {
			if actual.StorageAccountId.SubscriptionId != v.Expected.StorageAccountId.SubscriptionId {
				t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.StorageAccountId.SubscriptionId, actual.StorageAccountId.SubscriptionId)
			}
			if actual.StorageAccountId.ResourceGroup != v.Expected.StorageAccountId.ResourceGroup {
				t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.StorageAccountId.ResourceGroup, actual.StorageAccountId.ResourceGroup)
			}
			if actual.StorageAccountId.Name != v.Expected.StorageAccountId.Name {
				t.Fatalf("Expected %q but got %q for Name", v.Expected.StorageAccountId.Name, actual.StorageAccountId.Name)
			}
		}

		if actual.IPRule != nil && v.Expected.IPRule != nil && !strings.EqualFold(*actual.IPRule.IPAddressOrRange, *v.Expected.IPRule.IPAddressOrRange) {
			t.Fatalf("Expected %+v but got %+v for IPRule.IPAddressOrRange", *v.Expected.IPRule.IPAddressOrRange, *actual.IPRule.IPAddressOrRange)
		}

		if actual.VirtualNetworkRule != nil && v.Expected.VirtualNetworkRule != nil && !strings.EqualFold(*actual.VirtualNetworkRule.VirtualNetworkResourceID, *v.Expected.VirtualNetworkRule.VirtualNetworkResourceID) {
			t.Fatalf("Expected %+v but got %+v for VirtualNetworkRule.VirtualNetworkResourceID", *v.Expected.VirtualNetworkRule.VirtualNetworkResourceID, *actual.VirtualNetworkRule.VirtualNetworkResourceID)
		}

		if actual.ResourceAccessRule != nil && v.Expected.ResourceAccessRule != nil {
			if !strings.EqualFold(*actual.ResourceAccessRule.ResourceID, *v.Expected.ResourceAccessRule.ResourceID) {
				t.Fatalf("Expected %+v but got %+v for ResourceAccessRule.ResourceID", *v.Expected.ResourceAccessRule.ResourceID, *actual.ResourceAccessRule.ResourceID)
			}
			if !strings.EqualFold(*actual.ResourceAccessRule.TenantID, *v.Expected.ResourceAccessRule.TenantID) {
				t.Fatalf("Expected %+v but got %+v for ResourceAccessRule.TenantID", *v.Expected.ResourceAccessRule.TenantID, *actual.ResourceAccessRule.TenantID)
			}
		}
	}
}
