package parse

import "testing"

func TestMssqlDatabaseExtendedAuditingPolicy(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *MsSqlDatabaseExtendedAuditingPolicyId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Server Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Database",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1",
			Expected: nil,
		},
		{
			Name:     "Missing Sql Database Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/databases",
			Expected: nil,
		},
		{
			Name:     "Missing Extended Auditing Policy",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/databases/db1",
			Expected: nil,
		},
		{
			Name:     "Missing Extended Auditing Policy Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/databases/db1/extendedAuditingSettings",
			Expected: nil,
		},
		{
			Name:  "Extended Auditing Policy",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/databases/db1/extendedAuditingSettings/default",
			Expected: &MsSqlDatabaseExtendedAuditingPolicyId{
				ResourceGroup: "resGroup1",
				MsSqlServer:   "sqlServer1",
				MsDBName:      "db1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/databases/db1/ExtendedAuditingSettings/default",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := MssqlDatabaseExtendedAuditingPolicyID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.MsDBName != v.Expected.MsDBName {
			t.Fatalf("Expected %q but got %q for DB Name", v.Expected.MsDBName, actual.MsDBName)
		}

		if actual.MsSqlServer != v.Expected.MsSqlServer {
			t.Fatalf("Expected %q but got %q for Server Name", v.Expected.MsSqlServer, actual.MsSqlServer)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
