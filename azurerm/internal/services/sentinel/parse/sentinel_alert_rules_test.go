package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = SentinelAlertRuleId{}

func TestSentinelAlertRuleIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewSentinelAlertRuleID("group1", "space1", "rule1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

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
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "No Workspace ID",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights",
			Error: true,
		},
		{
			Name:  "No Alert Rule Name",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/",
			Error: true,
		},
		{
			Name:  "Incorrect Caseing",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/AlertRules/rule1",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1",
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
		if v.Error {
			t.Fatal("Expect an error but didn't get")
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

var _ resourceid.Formatter = SentinelAlertRuleActionId{}

func TestSentinelAlertRuleActionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewSentinelAlertRuleActionID("group1", "space1", "rule1", "action1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/action1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
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
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "No Workspace ID",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights",
			Error: true,
		},
		{
			Name:  "No Alert Rule Name",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/",
			Error: true,
		},
		{
			Name:  "No Alert Rule Action Name",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/",
			Error: true,
		},
		{
			Name:  "Incorrect Caseing",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/AlertRules/rule1/Actions/action1",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/space1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/action1",
			Expect: &SentinelAlertRuleActionId{
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
