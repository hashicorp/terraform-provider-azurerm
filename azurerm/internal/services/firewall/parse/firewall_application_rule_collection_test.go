package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = FirewallApplicationRuleCollectionId{}

func TestFirewallApplicationRuleCollectionIDFormatter(t *testing.T) {
	actual := NewFirewallApplicationRuleCollectionID("00000000-0000-0000-0000-000000000000", "mygroup1", "myfirewall", "applicationRuleCollection1").ID()
	expected := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/applicationRuleCollection1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFirewallApplicationRuleCollectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FirewallApplicationRuleCollectionId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},

		{
			// missing AzureFirewallName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for AzureFirewallName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/",
			Error: true,
		},

		{
			// missing ApplicationRuleCollectionName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/",
			Error: true,
		},

		{
			// missing value for ApplicationRuleCollectionName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/applicationRuleCollection1",
			Expected: &FirewallApplicationRuleCollectionId{
				SubscriptionId:                "00000000-0000-0000-0000-000000000000",
				ResourceGroup:                 "mygroup1",
				AzureFirewallName:             "myfirewall",
				ApplicationRuleCollectionName: "applicationRuleCollection1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/MYGROUP1/PROVIDERS/MICROSOFT.NETWORK/AZUREFIREWALLS/MYFIREWALL/APPLICATIONRULECOLLECTIONS/APPLICATIONRULECOLLECTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FirewallApplicationRuleCollectionID(v.Input)
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
		if actual.AzureFirewallName != v.Expected.AzureFirewallName {
			t.Fatalf("Expected %q but got %q for AzureFirewallName", v.Expected.AzureFirewallName, actual.AzureFirewallName)
		}
		if actual.ApplicationRuleCollectionName != v.Expected.ApplicationRuleCollectionName {
			t.Fatalf("Expected %q but got %q for ApplicationRuleCollectionName", v.Expected.ApplicationRuleCollectionName, actual.ApplicationRuleCollectionName)
		}
	}
}
