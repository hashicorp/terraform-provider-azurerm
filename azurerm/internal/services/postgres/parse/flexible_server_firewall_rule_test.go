package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = FlexibleServerFirewallRuleId{}

func TestFlexibleServerFirewallRuleIDFormatter(t *testing.T) {
	actual := NewFlexibleServerFirewallRuleID("12345678-1234-9876-4563-123456789012", "resGroup1", "flexibleServer1", "firewallRule1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/firewallRules/firewallRule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFlexibleServerFirewallRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FlexibleServerFirewallRuleId
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
			// missing FlexibleServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/",
			Error: true,
		},

		{
			// missing value for FlexibleServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/",
			Error: true,
		},

		{
			// missing FirewallRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/",
			Error: true,
		},

		{
			// missing value for FirewallRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/firewallRules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/firewallRules/firewallRule1",
			Expected: &FlexibleServerFirewallRuleId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				FlexibleServerName: "flexibleServer1",
				FirewallRuleName:   "firewallRule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DBFORPOSTGRESQL/FLEXIBLESERVERS/FLEXIBLESERVER1/FIREWALLRULES/FIREWALLRULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FlexibleServerFirewallRuleID(v.Input)
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
		if actual.FlexibleServerName != v.Expected.FlexibleServerName {
			t.Fatalf("Expected %q but got %q for FlexibleServerName", v.Expected.FlexibleServerName, actual.FlexibleServerName)
		}
		if actual.FirewallRuleName != v.Expected.FirewallRuleName {
			t.Fatalf("Expected %q but got %q for FirewallRuleName", v.Expected.FirewallRuleName, actual.FirewallRuleName)
		}
	}
}
