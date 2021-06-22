package parse

import (
	"testing"
)

func TestSynapseRoleAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *RoleAssignmentId
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
			Expected: &RoleAssignmentId{
				Scope:                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Synapse/workspaces/workspace1",
				DataPlaneAssignmentId: "00000000",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := RoleAssignmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for scope", v.Expected.Scope, actual.Scope)
		}

		if actual.DataPlaneAssignmentId != v.Expected.DataPlaneAssignmentId {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.DataPlaneAssignmentId, actual.DataPlaneAssignmentId)
		}
	}
}
