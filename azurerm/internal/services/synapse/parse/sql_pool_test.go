package parse

import (
	"testing"
)

func TestSynapseSqlPoolID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SqlPoolId
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
			Name:     "Missing SqlPool Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sqlPools",
			Expected: nil,
		},
		{
			Name:  "synapse SqlPool ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1/sqlPools/sqlPool1",
			Expected: &SqlPoolId{
				Workspace: &WorkspaceId{
					ResourceGroup: "resourceGroup1",
					Name:          "workspace1",
				},
				Name: "sqlPool1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1/SqlPools/sqlPool1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := SqlPoolID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Workspace.ResourceGroup != v.Expected.Workspace.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.Workspace.ResourceGroup, actual.Workspace.ResourceGroup)
		}

		if actual.Workspace.Name != v.Expected.Workspace.Name {
			t.Fatalf("Expected %q but got %q for Workspace Name", v.Expected.Workspace.Name, actual.Workspace.Name)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
