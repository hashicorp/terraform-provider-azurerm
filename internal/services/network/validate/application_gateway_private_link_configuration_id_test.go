package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestApplicationGatewayPrivateLinkConfigurationID(t *testing.T) {
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
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/",
			Valid: false,
		},

		{
			// missing ApplicationGatewayName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/",
			Valid: false,
		},

		{
			// missing value for ApplicationGatewayName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/",
			Valid: false,
		},

		{
			// missing PrivateLinkConfigurationName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/",
			Valid: false,
		},

		{
			// missing value for PrivateLinkConfigurationName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/privateLinkConfigurations/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/privateLinkConfigurations/hbtest",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00D16AA6-089C-400F-98EB-926ADB838EE1/RESOURCEGROUPS/CBUK-CORE-TESTOCE-APPGATEWAYSHARED-UKSOUTH/PROVIDERS/MICROSOFT.NETWORK/APPLICATIONGATEWAYS/CBUK-CORE-TESTOCE-APPGATEWAYSHARED-UKSOUTH/PRIVATELINKCONFIGURATIONS/HBTEST",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ApplicationGatewayPrivateLinkConfigurationID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
