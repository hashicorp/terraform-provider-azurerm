package parse

import (
	"testing"
)

func TestDedicatedHostID(t *testing.T) {
	testData := []struct {
		Name          string
		Input         string
		ExpectedOK    bool
		ExpectedValue *DedicatedHostId
	}{
		{
			Name:       "Empty",
			Input:      "",
			ExpectedOK: false,
		},
		{
			Name:       "No Resource Groups Segment",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectedOK: false,
		},
		{
			Name:       "No Resource Groups Value",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			ExpectedOK: false,
		},
		{
			Name:       "Resource Group ID",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			ExpectedOK: false,
		},
		{
			Name:       "Missing Host Group Value",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/",
			ExpectedOK: false,
		},
		{
			Name:       "Host Group ID",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/",
			ExpectedOK: false,
		},
		{
			Name:       "Missing Host Value",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/hosts/",
			ExpectedOK: false,
		},
		{
			Name:       "Host ID",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/hosts/host1",
			ExpectedOK: true,
			ExpectedValue: &DedicatedHostId{
				ResourceGroup: "resGroup1",
				HostGroup:     "group1",
				Name:          "host1",
			},
		},
		{
			Name:       "Wrong Casing",
			Input:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/group1/Hosts/host1",
			ExpectedOK: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := DedicatedHostID(v.Input)
		if err != nil {
			if v.ExpectedOK == false {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.ExpectedValue.Name {
			t.Fatalf("Expected %q but got %q for Name", v.ExpectedValue.Name, actual.Name)
		}

		if actual.HostGroup != v.ExpectedValue.HostGroup {
			t.Fatalf("Expected %q but got %q for HostGroup", v.ExpectedValue.HostGroup, actual.HostGroup)
		}

		if actual.ResourceGroup != v.ExpectedValue.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.ExpectedValue.ResourceGroup, actual.ResourceGroup)
		}
	}
}
