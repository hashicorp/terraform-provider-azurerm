package parse

import "testing"

func TestParseAppServiceEnvironmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AppServiceEnvironmentId
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
			Name:     "Missing environment name value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/hostingEnvironments/",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup1/providers/Microsoft.Web/hostingEnvironments/TestASEv2",
			Expected: &AppServiceEnvironmentId{
				ResourceGroup:          "testGroup1",
				HostingEnvironmentName: "TestASEv2",
			},
		},
		{
			Name:     "Wrong Case",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup1/providers/Microsoft.Web/HostingEnvironments/TestASEv2",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AppServiceEnvironmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.HostingEnvironmentName != v.Expected.HostingEnvironmentName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.HostingEnvironmentName, actual.HostingEnvironmentName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
