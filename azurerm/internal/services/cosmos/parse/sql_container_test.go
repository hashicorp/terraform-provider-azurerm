package parse

import (
	"testing"
)

func TestSqlContainerID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SqlContainerId
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
			Name:  "Missing SQL Container Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/",
			Error: true,
		},
		{
			Name:  "SQL Container ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/containers/container1",
			Error: false,
			Expect: &SqlContainerId{
				ResourceGroup:       "resGroup1",
				DatabaseAccountName: "acc1",
				SqlDatabaseName:     "db1",
				ContainerName:       "container1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/sqlDatabases/db1/Containers/container1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SqlContainerID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ContainerName != v.Expect.ContainerName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.ContainerName, actual.ContainerName)
		}

		if actual.SqlDatabaseName != v.Expect.SqlDatabaseName {
			t.Fatalf("Expected %q but got %q for Database", v.Expect.SqlDatabaseName, actual.SqlDatabaseName)
		}

		if actual.DatabaseAccountName != v.Expect.DatabaseAccountName {
			t.Fatalf("Expected %q but got %q for Account", v.Expect.DatabaseAccountName, actual.DatabaseAccountName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
