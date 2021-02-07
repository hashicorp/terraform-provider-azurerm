package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = CassandraTableId{}

func TestCassandraTableIDFormatter(t *testing.T) {
	actual := NewCassandraTableID("12345678-1234-9876-4563-123456789012", "resGroup1", "acc1", "keyspace1", "table1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/cassandraKeyspaces/keyspace1/tables/table1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestCassandraTableID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CassandraTableId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing DatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/",
			Error: true,
		},

		{
			// missing value for DatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/",
			Error: true,
		},

		{
			// missing CassandraKeyspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/",
			Error: true,
		},

		{
			// missing value for CassandraKeyspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/cassandraKeyspaces/",
			Error: true,
		},

		{
			// missing TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/cassandraKeyspaces/keyspace1/",
			Error: true,
		},

		{
			// missing value for TableName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/cassandraKeyspaces/keyspace1/tables/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/databaseAccounts/acc1/cassandraKeyspaces/keyspace1/tables/table1",
			Expected: &CassandraTableId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				DatabaseAccountName:   "acc1",
				CassandraKeyspaceName: "keyspace1",
				TableName:             "table1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DOCUMENTDB/DATABASEACCOUNTS/ACC1/CASSANDRAKEYSPACES/KEYSPACE1/TABLES/TABLE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CassandraTableID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.DatabaseAccountName != v.Expected.DatabaseAccountName {
			t.Fatalf("Expected %q but got %q for DatabaseAccountName", v.Expected.DatabaseAccountName, actual.DatabaseAccountName)
		}
		if actual.CassandraKeyspaceName != v.Expected.CassandraKeyspaceName {
			t.Fatalf("Expected %q but got %q for CassandraKeyspaceName", v.Expected.CassandraKeyspaceName, actual.CassandraKeyspaceName)
		}
		if actual.TableName != v.Expected.TableName {
			t.Fatalf("Expected %q but got %q for TableName", v.Expected.TableName, actual.TableName)
		}
	}
}
