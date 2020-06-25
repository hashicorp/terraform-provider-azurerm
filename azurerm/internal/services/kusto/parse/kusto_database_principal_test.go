package parse

import (
	"testing"
)

func TestKustoDatabasePrincipalId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *KustoDatabasePrincipalId
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
			Expected: &KustoDatabasePrincipalId{
				Name:          "aaduser=00000000-0000-0000-0000-000000000000;00000000-0000-0000-0000-000000000000",
				Role:          "Viewer",
				Database:      "database1",
				Cluster:       "cluster1",
				ResourceGroup: "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := KustoDatabasePrincipalID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.Role != v.Expected.Role {
			t.Fatalf("Expected %q but got %q for Role", v.Expected.Role, actual.Role)
		}

		if actual.Database != v.Expected.Database {
			t.Fatalf("Expected %q but got %q for Database", v.Expected.Database, actual.Database)
		}

		if actual.Cluster != v.Expected.Cluster {
			t.Fatalf("Expected %q but got %q for Cluster", v.Expected.Cluster, actual.Cluster)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
