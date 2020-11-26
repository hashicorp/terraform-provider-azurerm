package parse

import (
	"testing"
)

func TestKustoClusterPrincipalAssignmentId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ClusterPrincipalAssignmentId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Cluster",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/PrincipalAssignments/assignment1",
			Expected: nil,
		},
		{
			Name:     "Missing PrincipalAssignment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/",
			Expected: nil,
		},
		{
			Name:  "Cluster Principal Assignment ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/PrincipalAssignments/assignment1",
			Expected: &ClusterPrincipalAssignmentId{
				PrincipalAssignmentName: "assignment1",
				ClusterName:             "cluster1",
				ResourceGroup:           "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ClusterPrincipalAssignmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.PrincipalAssignmentName != v.Expected.PrincipalAssignmentName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.PrincipalAssignmentName, actual.PrincipalAssignmentName)
		}

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for Cluster", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
