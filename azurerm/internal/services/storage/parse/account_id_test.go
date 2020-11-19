package parse

import (
	"testing"
)

func TestParseAccountID(t *testing.T) {
	testData := []struct {
		input    string
		expected *AccountId
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello",
			expected: nil,
		},
		{
			// wrong case
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/storageaccounts/account1",
			expected: nil,
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/storageAccounts/account1",
			expected: &AccountId{
				Name:          "account1",
				ResourceGroup: "hello",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)
		actual, err := AccountID(v.input)

		// if we get something there shouldn't be an error
		if v.expected != nil && err == nil {
			continue
		}

		// if nothing's expected we should get an error
		if v.expected == nil && err != nil {
			continue
		}

		if v.expected == nil && actual == nil {
			continue
		}

		if v.expected == nil && actual != nil {
			t.Fatalf("Expected nothing but got %+v", actual)
		}
		if v.expected != nil && actual == nil {
			t.Fatalf("Expected %+v but got nil", actual)
		}

		if v.expected.ResourceGroup != actual.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but got %q", v.expected.ResourceGroup, actual.ResourceGroup)
		}
		if v.expected.Name != actual.Name {
			t.Fatalf("Expected Name to be %q but got %q", v.expected.Name, actual.Name)
		}
	}
}

func TestParseStorageSyncID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *StorageSyncId
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
			Name:     "Missing Storage Sync Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/",
			Expected: nil,
		},
		{
			Name:  "Storage Sync ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/Sync1",
			Expected: &StorageSyncId{
				Name:          "Sync1",
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/StorageSyncServices/Store1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseStorageSyncID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

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
				Name:            "group1",
				StorageSyncName: "sync1",
				ResourceGroup:   "resGroup1",
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

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.StorageSyncName != v.Expected.StorageSyncName {
			t.Fatalf("Expected %q but got %q for Storage Sync Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
