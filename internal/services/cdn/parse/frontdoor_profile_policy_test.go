package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FrontdoorProfilePolicyId{}

func TestFrontdoorProfilePolicyIDFormatter(t *testing.T) {
	actual := NewFrontdoorProfilePolicyID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "policy1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/policy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFrontdoorProfilePolicyID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontdoorProfilePolicyId
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
			// missing CdnWebApplicationFirewallPolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/",
			Error: true,
		},

		{
			// missing value for CdnWebApplicationFirewallPolicyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/cdnWebApplicationFirewallPolicies/policy1",
			Expected: &FrontdoorProfilePolicyId{
				SubscriptionId:                      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                       "resourceGroup1",
				CdnWebApplicationFirewallPolicyName: "policy1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CDN/CDNWEBAPPLICATIONFIREWALLPOLICIES/POLICY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FrontdoorProfilePolicyID(v.Input)
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
		if actual.CdnWebApplicationFirewallPolicyName != v.Expected.CdnWebApplicationFirewallPolicyName {
			t.Fatalf("Expected %q but got %q for CdnWebApplicationFirewallPolicyName", v.Expected.CdnWebApplicationFirewallPolicyName, actual.CdnWebApplicationFirewallPolicyName)
		}
	}
}
