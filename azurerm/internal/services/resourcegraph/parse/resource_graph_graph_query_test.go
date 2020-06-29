package parse

import (
	"testing"
)

func TestResourceGraphGraphQueryID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ResourceGraphGraphQueryId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Test Name",
			Input:    "your data",
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
			Name:     "Missing Resource Graph Query Name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ResourceGraph/queries/",
			Expected: nil,
		},
		{
			Name:  "Resource Graph Query Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ResourceGraph/queries/resGraph1",
			Expected: &ResourceGraphGraphQueryId{
				ResourceGroup: "resGroup1",
				ResourceName:  "resGraph1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ResourceGraph/Queries/resGraph1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := ResourceGraphGraphQueryID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.ResourceName != v.Expected.ResourceName {
			t.Fatalf("Expected %q but got %q for ResourceName", v.Expected.ResourceName, actual.ResourceName)
		}
	}
}
