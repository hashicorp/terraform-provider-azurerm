package parse

import (
	"testing"
)

func TestKustoDatabasePrincipalId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DatabasePrincipalId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing FQN",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Databases/database1",
			Expected: nil,
		},
		{
			Name:     "Missing Role",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/FQN/aaduser=;",
			Expected: nil,
		},
		{
			Name:  "Database Principal ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/Role/Viewer/FQN/aaduser=00000000-0000-0000-0000-000000000000;00000000-0000-0000-0000-000000000000",
			Expected: &DatabasePrincipalId{
				FQNName:       "aaduser=00000000-0000-0000-0000-000000000000;00000000-0000-0000-0000-000000000000",
				RoleName:      "Viewer",
				DatabaseName:  "database1",
				ClusterName:   "cluster1",
				ResourceGroup: "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := DatabasePrincipalID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.FQNName != v.Expected.FQNName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.FQNName, actual.FQNName)
		}

		if actual.RoleName != v.Expected.RoleName {
			t.Fatalf("Expected %q but got %q for Role", v.Expected.RoleName, actual.RoleName)
		}

		if actual.DatabaseName != v.Expected.DatabaseName {
			t.Fatalf("Expected %q but got %q for Database", v.Expected.DatabaseName, actual.DatabaseName)
		}

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for Cluster", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
