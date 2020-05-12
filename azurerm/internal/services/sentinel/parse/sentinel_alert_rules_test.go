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
				Subscription:  "00000000-0000-0000-0000-000000000000",
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

		if actual.Subscription != v.Expect.Subscription {
			t.Fatalf("Expected %q but got %q for Subscription", v.Expect.Subscription, actual.Subscription)
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

func TestSentinelAlertRuleActionID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SentinelAlertRuleActionId
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
			Name:  "No Alert Rule Action Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/",
			Error: true,
		},
		{
			Name:  "Incorrect Caseing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/AlertRules/rule1/Actions/action1",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/action1",
			Expect: &SentinelAlertRuleActionId{
				Subscription:  "00000000-0000-0000-0000-000000000000",
				ResourceGroup: "resGroup1",
				Workspace:     "space1",
				Rule:          "rule1",
				Name:          "action1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SentinelAlertRuleActionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get")
		}

		if actual.Subscription != v.Expect.Subscription {
			t.Fatalf("Expected %q but got %q for Subscription", v.Expect.Subscription, actual.Subscription)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Workspace != v.Expect.Workspace {
			t.Fatalf("Expected %q but got %q for Workspace", v.Expect.Name, actual.Name)
		}

		if actual.Rule != v.Expect.Rule {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Rule, actual.Rule)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
