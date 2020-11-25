package parse

import (
	"testing"
)

func TestMongoDbCollectionId(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *MongoDbCollectionId
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
			Name:  "Missing MongoDB Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/mongodbDatabases/",
			Error: true,
		},
		{
			Name:  "Missing Collection Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/mongodbDatabases/db1/collections/",
			Error: true,
		},
		{
			Name:  "MongoDB Collection ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/mongodbDatabases/db1/collections/coll1",
			Error: false,
			Expect: &MongoDbCollectionId{
				ResourceGroup:       "resGroup1",
				DatabaseAccountName: "acc1",
				MongodbDatabaseName: "db1",
				CollectionName:      "coll1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/MongodbDatabases/db1/Collections/coll1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := MongoDbCollectionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.CollectionName != v.Expect.CollectionName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.CollectionName, actual.CollectionName)
		}

		if actual.MongodbDatabaseName != v.Expect.MongodbDatabaseName {
			t.Fatalf("Expected %q but got %q for Database", v.Expect.MongodbDatabaseName, actual.MongodbDatabaseName)
		}

		if actual.DatabaseAccountName != v.Expect.DatabaseAccountName {
			t.Fatalf("Expected %q but got %q for Account", v.Expect.DatabaseAccountName, actual.DatabaseAccountName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
