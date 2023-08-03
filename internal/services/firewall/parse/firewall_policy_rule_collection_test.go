package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FirewallPolicyRuleCollectionId{}

func TestFirewallPolicyRuleCollectionIDFormatter(t *testing.T) {
	actual := NewFirewallPolicyRuleCollectionID("12345678-1234-9876-4563-123456789012", "resGroup1", "policy1", "ruleCollectionGroup1", "ruleCollection1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1/ruleCollections/ruleCollection1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFirewallPolicyRuleCollectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FirewallPolicyRuleCollectionId
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
			// missing FirewallPolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FirewallPolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/",
			Error: true,
		},

		{
			// missing RuleCollectionGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/",
			Error: true,
		},

		{
			// missing value for RuleCollectionGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/",
			Error: true,
		},

		{
			// missing RuleCollectionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1/",
			Error: true,
		},

		{
			// missing value for RuleCollectionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1/ruleCollections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1/ruleCollections/ruleCollection1",
			Expected: &FirewallPolicyRuleCollectionId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "resGroup1",
				FirewallPolicyName:      "policy1",
				RuleCollectionGroupName: "ruleCollectionGroup1",
				RuleCollectionName:      "ruleCollection1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/FIREWALLPOLICIES/POLICY1/RULECOLLECTIONGROUPS/RULECOLLECTIONGROUP1/RULECOLLECTIONS/RULECOLLECTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FirewallPolicyRuleCollectionID(v.Input)
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
		if actual.FirewallPolicyName != v.Expected.FirewallPolicyName {
			t.Fatalf("Expected %q but got %q for FirewallPolicyName", v.Expected.FirewallPolicyName, actual.FirewallPolicyName)
		}
		if actual.RuleCollectionGroupName != v.Expected.RuleCollectionGroupName {
			t.Fatalf("Expected %q but got %q for RuleCollectionGroupName", v.Expected.RuleCollectionGroupName, actual.RuleCollectionGroupName)
		}
		if actual.RuleCollectionName != v.Expected.RuleCollectionName {
			t.Fatalf("Expected %q but got %q for RuleCollectionName", v.Expected.RuleCollectionName, actual.RuleCollectionName)
		}
	}
}
