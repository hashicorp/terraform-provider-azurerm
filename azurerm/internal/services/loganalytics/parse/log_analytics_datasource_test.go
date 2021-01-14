package parse

import (
	"testing"
)

func TestLogAnalyticsDataSourceID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *LogAnalyticsDataSourceId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/",
			Error: true,
		},
		{
			Name:  "Missing Workspace Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces",
			Error: true,
		},
		{
			Name:  "Workspace Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1",
			Error: true,
		},
		{
			Name:  "Missing DataSource Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources",
			Error: true,
		},
		{
			Name:  "DataSource Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/datasources/datasource1",
			Error: true,
			Expect: &LogAnalyticsDataSourceId{
				ResourceGroup: "resGroup1",
				Workspace:     "workspace1",
				Name:          "datasource1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := LogAnalyticsDataSourceID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Workspace != v.Expect.Workspace {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Workspace, actual.Workspace)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
