package webapplicationfirewallpolicies

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FrontDoorWebApplicationFirewallPoliciesId{}

func TestFrontDoorWebApplicationFirewallPoliciesIDFormatter(t *testing.T) {
	actual := NewFrontDoorWebApplicationFirewallPoliciesID("{subscriptionId}", "{resourceGroupName}", "{policyName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{policyName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseFrontDoorWebApplicationFirewallPoliciesID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontDoorWebApplicationFirewallPoliciesId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing FrontDoorWebApplicationFirewallPolicyName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FrontDoorWebApplicationFirewallPolicyName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{policyName}",
			Expected: &FrontDoorWebApplicationFirewallPoliciesId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				FrontDoorWebApplicationFirewallPolicyName: "{policyName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.NETWORK/FRONTDOORWEBAPPLICATIONFIREWALLPOLICIES/{POLICYNAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFrontDoorWebApplicationFirewallPoliciesID(v.Input)
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
		if actual.FrontDoorWebApplicationFirewallPolicyName != v.Expected.FrontDoorWebApplicationFirewallPolicyName {
			t.Fatalf("Expected %q but got %q for FrontDoorWebApplicationFirewallPolicyName", v.Expected.FrontDoorWebApplicationFirewallPolicyName, actual.FrontDoorWebApplicationFirewallPolicyName)
		}
	}
}

func TestParseFrontDoorWebApplicationFirewallPoliciesIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontDoorWebApplicationFirewallPoliciesId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing FrontDoorWebApplicationFirewallPolicyName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FrontDoorWebApplicationFirewallPolicyName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/{policyName}",
			Expected: &FrontDoorWebApplicationFirewallPoliciesId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				FrontDoorWebApplicationFirewallPolicyName: "{policyName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontdoorwebapplicationfirewallpolicies/{policyName}",
			Expected: &FrontDoorWebApplicationFirewallPoliciesId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				FrontDoorWebApplicationFirewallPolicyName: "{policyName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/FRONTDOORWEBAPPLICATIONFIREWALLPOLICIES/{policyName}",
			Expected: &FrontDoorWebApplicationFirewallPoliciesId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				FrontDoorWebApplicationFirewallPolicyName: "{policyName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/FrOnTdOoRwEbApPlIcAtIoNfIrEwAlLpOlIcIeS/{policyName}",
			Expected: &FrontDoorWebApplicationFirewallPoliciesId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				FrontDoorWebApplicationFirewallPolicyName: "{policyName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFrontDoorWebApplicationFirewallPoliciesIDInsensitively(v.Input)
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
		if actual.FrontDoorWebApplicationFirewallPolicyName != v.Expected.FrontDoorWebApplicationFirewallPolicyName {
			t.Fatalf("Expected %q but got %q for FrontDoorWebApplicationFirewallPolicyName", v.Expected.FrontDoorWebApplicationFirewallPolicyName, actual.FrontDoorWebApplicationFirewallPolicyName)
		}
	}
}
