package parse

import (
	"testing"
)

func TestLogAnalyticsLinkedServiceID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *LogAnalyticsLinkedServiceId
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
			Name:  "Resource Group Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedServices/Cluster",
			Expected: &LogAnalyticsLinkedServiceId{
				ResourceGroup: "resourceGroup1",
				WorkspaceName: "workspace1",
				Type:          "Cluster",
			},
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing LinkedService Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedServices",
			Expected: nil,
		},
		{
			Name:  "Log Analytics Linked Service ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedServices/Automation",
			Expected: &LogAnalyticsLinkedServiceId{
				ResourceGroup: "resourceGroup1",
				WorkspaceName: "workspace1",
				Type:          "Automation",
			},
		},
		{
			Name:     "LinkedServices Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/LinkedServices/Automation",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := LogAnalyticsLinkedServiceID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}

		if actual.Type != v.Expected.Type {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Type, actual.Type)
		}
	}
}
