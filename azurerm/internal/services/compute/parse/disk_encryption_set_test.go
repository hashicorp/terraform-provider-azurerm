package parse

import "testing"

func TestDiskEncryptionSetID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *DiskEncryptionSetId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Group Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing Disk Encryption Set Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/diskEncryptionSets/",
			Error: true,
		},
		{
			Name:  "Disk Encryption Set ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/diskEncryptionSets/encryptionSet1",
			Error: false,
			Expect: &DiskEncryptionSetId{
				ResourceGroup: "resGroup1",
				Name:          "encryptionSet1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/DiskEncryptionSets/encryptionSet1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := DiskEncryptionSetID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q", v.Expect.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
