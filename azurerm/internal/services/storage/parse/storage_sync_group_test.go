package parse

import "testing"

func TestParseStorageSyncGroupID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *StorageSyncGroupId
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
			Name:     "Missing Sync Group",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/sync1",
			Expected: nil,
		},
		{
			Name:     "Missing Sync Group Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/sync1/syncGroups/",
			Expected: nil,
		},
		{
			Name:  "Sync Group Id",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/sync1/syncGroups/group1",
			Expected: &StorageSyncGroupId{
				SyncGroupName:          "group1",
				StorageSyncServiceName: "sync1",
				ResourceGroup:          "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/sync1/SyncGroups/group1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := StorageSyncGroupID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SyncGroupName != v.Expected.SyncGroupName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.SyncGroupName, actual.SyncGroupName)
		}

		if actual.StorageSyncServiceName != v.Expected.StorageSyncServiceName {
			t.Fatalf("Expected %q but got %q for Storage Sync Name", v.Expected.SyncGroupName, actual.SyncGroupName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
