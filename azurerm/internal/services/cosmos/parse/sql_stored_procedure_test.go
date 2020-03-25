package parse

import (
	"testing"
)

func TestStoredProcedureID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *StoredProcedureId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing Database Account Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/",
			Error: true,
		},
		{
			Name:  "Missing SQL Database Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/",
			Error: true,
		},
		{
			Name:  "Missing Container Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/",
			Error: true,
		},
		{
			Name:  "Missing Stored Procedure Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/container1/storedProcedures/",
			Error: true,
		},
		{
			Name:  "Stored Procedure ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/container1/storedProcedures/sproc1",
			Error: false,
			Expect: &StoredProcedureId{
				ResourceGroup: "resGroup1",
				Account:       "acc1",
				Database:      "db1",
				Container:     "container1",
				Name:          "sproc1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/container1/StoredProcedures/sproc1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := StoredProcedureID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.Container != v.Expect.Container {
			t.Fatalf("Expected %q but got %q for Container", v.Expect.Container, actual.Container)
		}

		if actual.Database != v.Expect.Database {
			t.Fatalf("Expected %q but got %q for Database", v.Expect.Database, actual.Database)
		}

		if actual.Account != v.Expect.Account {
			t.Fatalf("Expected %q but got %q for Account", v.Expect.Account, actual.Account)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
