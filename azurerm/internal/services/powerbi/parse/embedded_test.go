package parse

import (
	"testing"
)

func TestPowerBIEmbeddedId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *EmbeddedId
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
			Name:     "Missing PowerBI Embedded value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.PowerBIDedicated/capacities",
			Expected: nil,
		},
		{
			Name:  "PowerBI Embedded ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.PowerBIDedicated/capacities/capacity1",
			Expected: &EmbeddedId{
				ResourceGroup: "resGroup1",
				CapacityName:  "capacity1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := EmbeddedID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.CapacityName != v.Expected.CapacityName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.CapacityName, actual.CapacityName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
