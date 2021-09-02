package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RulesEngineId{}

func TestRulesEngineIDFormatter(t *testing.T) {
	actual := NewRulesEngineID("c45eeda7-1811-4ab1-8fe2-efdd99c9d489", "FrontDoorExampleResourceGroup", "heoelri-example-fd", "rule1").ID()
	expected := "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/rulesengines/rule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRulesEngineID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RulesEngineId
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
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/",
			Error: true,
		},

		{
			// missing FrontdoorName
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FrontdoorName
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/rulesengines/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/rulesengines/rule1",
			Expected: &RulesEngineId{
				SubscriptionId: "c45eeda7-1811-4ab1-8fe2-efdd99c9d489",
				ResourceGroup:  "FrontDoorExampleResourceGroup",
				FrontdoorName:  "heoelri-example-fd",
				Name:           "rule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/C45EEDA7-1811-4AB1-8FE2-EFDD99C9D489/RESOURCEGROUPS/FRONTDOOREXAMPLERESOURCEGROUP/PROVIDERS/MICROSOFT.NETWORK/FRONTDOORS/HEOELRI-EXAMPLE-FD/RULESENGINES/RULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RulesEngineID(v.Input)
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
		if actual.FrontdoorName != v.Expected.FrontdoorName {
			t.Fatalf("Expected %q but got %q for FrontdoorName", v.Expected.FrontdoorName, actual.FrontdoorName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
