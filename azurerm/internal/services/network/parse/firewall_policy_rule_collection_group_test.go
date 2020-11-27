package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = FirewallPolicyRuleCollectionGroupId{}

func TestFirewallPolicyRuleCollectionGroupIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewFirewallPolicyRuleCollectionGroupID(NewFirewallPolicyID("group1", "policy1"), "group1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/group1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFirewallPolicyRuleCollectionGroupID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *FirewallPolicyRuleCollectionGroupId
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
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/",
			Error: true,
		},
		{
			Name:  "No Policy Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies",
			Error: true,
		},
		{
			Name:  "No Group Name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups",
			Error: true,
		},
		{
			Name:  "The ID of old API version (2020-03-01)",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleGroups/rulegroup1",
			Error: true,
		},
		{
			Name:  "Correct Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/rulegroup1",
			Expect: &FirewallPolicyRuleCollectionGroupId{
				ResourceGroup: "group1",
				PolicyName:    "policy1",
				Name:          "rulegroup1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := FirewallPolicyRuleCollectionGroupID(v.Input)
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

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
