package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = LogAnalyticsLinkedStorageAccountId{}

func TestLogAnalyticsLinkedStorageAccountIDFormatter(t *testing.T) {
	actual := NewLogAnalyticsLinkedStorageAccountID("12345678-1234-9876-4563-123456789012", "resGroup1", "workspace1", "query").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/workspaces/workspace1/linkedStorageAccounts/query"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestLogAnalyticsLinkedStorageAccountID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *LogAnalyticsLinkedStorageAccountId
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
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/workspaces/",
			Error: true,
		},

		{
			// missing LinkedStorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for LinkedStorageAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/workspaces/workspace1/linkedStorageAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/microsoft.operationalinsights/workspaces/workspace1/linkedStorageAccounts/query",
			Expected: &LogAnalyticsLinkedStorageAccountId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resGroup1",
				WorkspaceName:            "workspace1",
				LinkedStorageAccountName: "query",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/WORKSPACE1/LINKEDSTORAGEACCOUNTS/QUERY",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := LogAnalyticsLinkedStorageAccountID(v.Input)
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
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.LinkedStorageAccountName != v.Expected.LinkedStorageAccountName {
			t.Fatalf("Expected %q but got %q for LinkedStorageAccountName", v.Expected.LinkedStorageAccountName, actual.LinkedStorageAccountName)
		}
	}
}
