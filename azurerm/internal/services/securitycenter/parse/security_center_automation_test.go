package parse

import "testing"

func TestSecurityCentreAutomationID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SecurityCenterAutomationId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},
		{
			Name:  "No Automation Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG1/providers/Microsoft.Security/",
			Error: true,
		},
		{
			Name:  "No Automation Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG1/providers/Microsoft.Security/automations/",
			Error: true,
		},
		{
			Name:  "Security Center Automation ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG1/providers/Microsoft.Security/automations/testAutomation",
			Expect: &SecurityCenterAutomationId{
				ResourceGroup:  "testRG1",
				AutomationName: "testAutomation",
			},
		},
		{
			Name:  "Wrong Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG1/providers/Microsoft.Security/Automations/testAutomation",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SecurityCenterAutomationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.AutomationName != v.Expect.AutomationName {
			t.Fatalf("Expected %q but got %q for Automation Name", v.Expect.AutomationName, actual.AutomationName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group Name", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
