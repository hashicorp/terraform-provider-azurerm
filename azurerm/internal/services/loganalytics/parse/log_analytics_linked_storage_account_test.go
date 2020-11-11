package parse

import (
	"testing"
)

func TestOperationalinsightsLinkedStorageAccountID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *LogAnalyticsLinkedStorageAccountId
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
			Name:     "Missing LinkedStorageAccount Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedStorageAccounts",
			Expected: nil,
		},
		{
			Name:  "Log Analytics Linked Storage Account ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedStorageAccounts/dataSourceType1",
			Expected: &LogAnalyticsLinkedStorageAccountId{
				ResourceGroup: "resourceGroup1",
				WorkspaceName: "workspace1",
				WorkspaceID:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1",
				Name:          "dataSourceType1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/LinkedStorageAccounts/dataSourceType1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := LogAnalyticsLinkedStorageAccountID(v.Input)
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

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
