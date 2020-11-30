package parse

import "testing"

func TestMsSqlRestoreDBID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *RestorableDroppedDatabaseId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Restore Name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000,000000000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/,000000000000000000",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/,000000000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Server Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/,000000000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Restorable Database",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1,000000000000000000",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Restorable Database Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/restorableDroppedDatabases,000000000000000000",
			Expected: nil,
		},
		{
			Name:  "Sql Restorable Database ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/restorableDroppedDatabases/sqlDB1,000000000000000000",
			Expected: &RestorableDroppedDatabaseId{
				Name:          "sqlDB1",
				MsSqlServer:   "sqlServer1",
				ResourceGroup: "resGroup1",
				RestoreName:   "000000000000000000",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/RestorableDroppedDatabases/sqlDB1,000000000000000000",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := RestorableDroppedDatabaseID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.RestoreName != v.Expected.RestoreName {
			t.Fatalf("Expected %q but got %q for Restore Name", v.Expected.Name, actual.Name)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.MsSqlServer != v.Expected.MsSqlServer {
			t.Fatalf("Expected %q but got %q for Sql Server", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
