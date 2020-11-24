package parse

import (
	"testing"
)

func TestAutomationConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ConnectionId
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
			Name:     "Missing Automation Accounts Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Automation/automationAccounts/",
			Expected: nil,
		},
		{
			Name:  "Automation Connection ID",
			Input: "/subscriptions/85b3dbca-5974-4067-9669-67a141095a76/resourceGroups/resGroup1/providers/Microsoft.Automation/automationAccounts/account1/connections/conn1",
			Expected: &ConnectionId{
				AutomationAccountName: "account1",
				ConnectionName:        "conn1",
				ResourceGroup:         "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/85b3dbca-5974-4067-9669-67a141095a76/resourceGroups/resGroup1/providers/Microsoft.Automation/automationAccounts/account1/Connections/conn1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.AutomationAccountName != v.Expected.AutomationAccountName {
			t.Fatalf("Expected %q but got %q for AccountName", v.Expected.AutomationAccountName, actual.AutomationAccountName)
		}

		if actual.ConnectionName != v.Expected.ConnectionName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.ConnectionName, actual.ConnectionName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
