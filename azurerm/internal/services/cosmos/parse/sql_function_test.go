package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SqlFunctionId{}

func TestSqlFunctionIDFormatter(t *testing.T) {
	actual := NewSqlFunctionID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "account1", "database1", "container1", "userDefinedFunction1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/container1/userDefinedFunctions/userDefinedFunction1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSqlFunctionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SqlFunctionId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/",
			Error: true,
		},

		{
			// missing value for DatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/",
			Error: true,
		},

		{
			// missing SqlDatabaseName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/",
			Error: true,
		},

		{
			// missing value for SqlDatabaseName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/",
			Error: true,
		},

		{
			// missing ContainerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/",
			Error: true,
		},

		{
			// missing value for ContainerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/",
			Error: true,
		},

		{
			// missing UserDefinedFunctionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/container1/",
			Error: true,
		},

		{
			// missing value for UserDefinedFunctionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/container1/userDefinedFunctions/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.DocumentDB/databaseAccounts/account1/sqlDatabases/database1/containers/container1/userDefinedFunctions/userDefinedFunction1",
			Expected: &SqlFunctionId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "resourceGroup1",
				DatabaseAccountName:     "account1",
				SqlDatabaseName:         "database1",
				ContainerName:           "container1",
				UserDefinedFunctionName: "userDefinedFunction1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.DOCUMENTDB/DATABASEACCOUNTS/ACCOUNT1/SQLDATABASES/DATABASE1/CONTAINERS/CONTAINER1/USERDEFINEDFUNCTIONS/USERDEFINEDFUNCTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SqlFunctionID(v.Input)
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
		if actual.SqlDatabaseName != v.Expected.SqlDatabaseName {
			t.Fatalf("Expected %q but got %q for SqlDatabaseName", v.Expected.SqlDatabaseName, actual.SqlDatabaseName)
		}
		if actual.ContainerName != v.Expected.ContainerName {
			t.Fatalf("Expected %q but got %q for ContainerName", v.Expected.ContainerName, actual.ContainerName)
		}
		if actual.UserDefinedFunctionName != v.Expected.UserDefinedFunctionName {
			t.Fatalf("Expected %q but got %q for UserDefinedFunctionName", v.Expected.UserDefinedFunctionName, actual.UserDefinedFunctionName)
		}
	}
}
