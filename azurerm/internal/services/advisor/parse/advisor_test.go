package parse

import (
	"testing"
)

func TestAdvisorSubscriptionID(t *testing.T) {
	testData := []struct {
		Name        string
		Input       string
		ExpectError bool
	}{
		{
			Name:        "Empty",
			Input:       "",
			ExpectError: true,
		},
		{
			Name:        "No Configurations Segment",
			Input:       "/subscriptions/00000000-0000-0000-0000-000000000000",
			ExpectError: true,
		},
		{
			Name:        "Have Resource Group ID",
			Input:       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			ExpectError: true,
		},
		{
			Name:        "Right ID",
			Input:       "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Advisor/configurations/nameId",
			ExpectError: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		err := AdvisorSubscriptionID(v.Input)
		if err != nil {
			if v.ExpectError {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}
	}
}

func TestAdvisorResGroupID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AdvisorResGroupId
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
			Name:     "Missing Name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Advisor/configurations/",
			Expected: nil,
		},
		{
			Name:  "App Configuration ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Advisor/configurations/name1",
			Expected: &AdvisorResGroupId{
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Advisor/Configurations/name1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AdvisorResGroupID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
