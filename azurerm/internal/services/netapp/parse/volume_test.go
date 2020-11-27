package parse

import (
	"testing"
)

func TestNetAppVolumeId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VolumeId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing NetApp Account Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/",
			Expected: nil,
		},
		{
			Name:     "NetApp Account ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1",
			Expected: nil,
		},
		{
			Name:     "Missing NetApp Pool Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/",
			Expected: nil,
		},
		{
			Name:     "NetApp Pool ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1",
			Expected: nil,
		},
		{
			Name:     "Missing NetApp Volume Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/",
			Expected: nil,
		},
		{
			Name:  "NetApp Volume ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1",
			Expected: &VolumeId{
				Name:              "volume1",
				CapacityPoolName:  "pool1",
				NetAppAccountName: "account1",
				ResourceGroup:     "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/VOLUMES/volume1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VolumeID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.NetAppAccountName != v.Expected.NetAppAccountName {
			t.Fatalf("Expected %q but got %q for Account Name", v.Expected.NetAppAccountName, actual.NetAppAccountName)
		}

		if actual.CapacityPoolName != v.Expected.CapacityPoolName {
			t.Fatalf("Expected %q but got %q for Pool Name", v.Expected.CapacityPoolName, actual.CapacityPoolName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
