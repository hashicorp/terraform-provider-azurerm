package parse

import (
	"testing"
)

func TestGremlinGraphID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *GremlinGraphId
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
			Name:  "Missing Gremlin Database Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/gremlinDatabases/",
			Error: true,
		},
		{
			Name:  "Missing Gremlin Graph Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/gremlinDatabases/graphs",
			Error: true,
		},
		{
			Name:  "Gremlin Graph ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/gremlinDatabases/database1/graphs/graph1",
			Error: false,
			Expect: &GremlinGraphId{
				ResourceGroup:       "resGroup1",
				DatabaseAccountName: "acc1",
				GremlinDatabaseName: "database1",
				GraphName:           "graph1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/gremlinDatabases/database1/Graphs/graph1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := GremlinGraphID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.GraphName != v.Expect.GraphName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.GraphName, actual.GraphName)
		}

		if actual.GremlinDatabaseName != v.Expect.GremlinDatabaseName {
			t.Fatalf("Expected %q but got %q for Database", v.Expect.GremlinDatabaseName, actual.GremlinDatabaseName)
		}

		if actual.DatabaseAccountName != v.Expect.DatabaseAccountName {
			t.Fatalf("Expected %q but got %q for Account", v.Expect.DatabaseAccountName, actual.DatabaseAccountName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
