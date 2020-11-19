package parse

import (
	"testing"
)

func TestParseStorageContainerResourceManagerID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *StorageContainerResourceManagerId
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
			Name:     "Missing Containers Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/",
			Expected: nil,
		},
		{
			Name:  "Storage Containers Resource Manager ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1",
			Expected: &StorageContainerResourceManagerId{
				ResourceGroup:   "resGroup1",
				AccountName:     "account1",
				BlobServiceName: "default",
				Name:            "container1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/Containers/container1",
			Expected: nil,
		},
		{
			Name:     "Storage Container Data Plane ID",
			Input:    "https://account1.blob.core.windows.net/container1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := StorageContainerResourceManagerID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.BlobServiceName != v.Expected.BlobServiceName {
			t.Fatalf("Expected %q but got %q for Blob Service Name", v.Expected.BlobServiceName, actual.BlobServiceName)
		}

		if actual.AccountName != v.Expected.AccountName {
			t.Fatalf("Expected %q but got %q for Account Name", v.Expected.AccountName, actual.AccountName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
