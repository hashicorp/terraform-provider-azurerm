package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SnapshotId{}

func TestSnapshotIDFormatter(t *testing.T) {
	actual := NewSnapshotID("12345678-1234-9876-4563-123456789012", "resGroup1", "account1", "pool1", "volume1", "snapshot1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/snapshot1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSnapshotID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SnapshotId
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
			// missing NetAppAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/",
			Error: true,
		},

		{
			// missing value for NetAppAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/",
			Error: true,
		},

		{
			// missing CapacityPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/",
			Error: true,
		},

		{
			// missing value for CapacityPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/",
			Error: true,
		},

		{
			// missing VolumeName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/",
			Error: true,
		},

		{
			// missing value for VolumeName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/snapshot1",
			Expected: &SnapshotId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				NetAppAccountName: "account1",
				CapacityPoolName:  "pool1",
				VolumeName:        "volume1",
				Name:              "snapshot1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETAPP/NETAPPACCOUNTS/ACCOUNT1/CAPACITYPOOLS/POOL1/VOLUMES/VOLUME1/SNAPSHOTS/SNAPSHOT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SnapshotID(v.Input)
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
		if actual.NetAppAccountName != v.Expected.NetAppAccountName {
			t.Fatalf("Expected %q but got %q for NetAppAccountName", v.Expected.NetAppAccountName, actual.NetAppAccountName)
		}
		if actual.CapacityPoolName != v.Expected.CapacityPoolName {
			t.Fatalf("Expected %q but got %q for CapacityPoolName", v.Expected.CapacityPoolName, actual.CapacityPoolName)
		}
		if actual.VolumeName != v.Expected.VolumeName {
			t.Fatalf("Expected %q but got %q for VolumeName", v.Expected.VolumeName, actual.VolumeName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
