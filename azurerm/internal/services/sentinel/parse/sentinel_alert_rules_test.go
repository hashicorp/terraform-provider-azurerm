package parse

import (
	"testing"
)

func TestSentinelAlertRuleID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SentinelAlertRuleId
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
			Name:  "No Workspace ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights",
			Error: true,
		},
		{
			Name:  "No Alert Rule Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/",
			Error: true,
		},
		{
			Name:  "Incorrect Provider for Alert Rules",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Foo.Bar/alertRules/rule1",
			Error: true,
		},
		{
			Name:  "Incorrect Caseing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/AlertRules/rule1",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1",
			Expect: &SentinelAlertRuleId{
				ResourceGroup: "resGroup1",
				Workspace:     "space1",
				Name:          "rule1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SentinelAlertRuleID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Workspace != v.Expect.Workspace {
			t.Fatalf("Expected %q but got %q for Workspace", v.Expect.Name, actual.Name)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
