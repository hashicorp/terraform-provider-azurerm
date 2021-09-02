package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestRulesEngineID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/",
			Valid: false,
		},

		{
			// missing FrontdoorName
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/",
			Valid: false,
		},

		{
			// missing value for FrontdoorName
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/",
			Valid: false,
		},

		{
			// missing Name
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/rulesengines/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/c45eeda7-1811-4ab1-8fe2-efdd99c9d489/resourceGroups/FrontDoorExampleResourceGroup/providers/Microsoft.Network/frontdoors/heoelri-example-fd/rulesengines/rule1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/C45EEDA7-1811-4AB1-8FE2-EFDD99C9D489/RESOURCEGROUPS/FRONTDOOREXAMPLERESOURCEGROUP/PROVIDERS/MICROSOFT.NETWORK/FRONTDOORS/HEOELRI-EXAMPLE-FD/RULESENGINES/RULE1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RulesEngineID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
