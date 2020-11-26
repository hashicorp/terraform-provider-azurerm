package parse

import "testing"

func TestMssqlDatabaseExtendedAuditingPolicy(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *DatabaseExtendedAuditingPolicyId
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
			Expected: &DatabaseExtendedAuditingPolicyId{
				ResourceGroup: "resGroup1",
				ServerName:    "sqlServer1",
				DatabaseName:  "db1",
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

		actual, err := DatabaseExtendedAuditingPolicyID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.DatabaseName != v.Expected.DatabaseName {
			t.Fatalf("Expected %q but got %q for DB Name", v.Expected.DatabaseName, actual.DatabaseName)
		}

		if actual.ServerName != v.Expected.ServerName {
			t.Fatalf("Expected %q but got %q for Server Name", v.Expected.ServerName, actual.ServerName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
