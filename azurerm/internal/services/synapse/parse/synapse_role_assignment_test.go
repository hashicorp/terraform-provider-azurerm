package parse

import (
	"testing"
)

func TestSynapseRoleAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SynapseRoleAssignmentId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Synapse Role Assignment ID part",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1",
			Expected: nil,
		},
		{
			Name:     "Missing Synapse Workspace ID part",
			Input:    "00000000",
			Expected: nil,
		},
		{
			Name:  "synapse Role Assignment ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1|00000000",
			Expected: &SynapseRoleAssignmentId{
				Workspace: &SynapseWorkspaceId{
					ResourceGroup: "resourceGroup1",
					Name:          "workspace1",
				},
				Id: "00000000",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/Workspaces/workspace1|00000000",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := SynapseRoleAssignmentID(v.Input)
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
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.Workspace.Name, actual.Workspace.Name)
		}

		if actual.Id != v.Expected.Id {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Id, actual.Id)
		}
	}
}
