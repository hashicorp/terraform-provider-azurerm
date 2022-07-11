package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ApplicationGatewayPrivateLinkConfigurationId{}

func TestApplicationGatewayPrivateLinkConfigurationIDFormatter(t *testing.T) {
	actual := NewApplicationGatewayPrivateLinkConfigurationID("00d16aa6-089c-400f-98eb-926adb838ee1", "cbuk-core-testoce-appgatewayshared-uksouth", "cbuk-core-testoce-appgatewayshared-uksouth", "hbtest").ID()
	expected := "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/privateLinkConfigurations/hbtest"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestApplicationGatewayPrivateLinkConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ApplicationGatewayPrivateLinkConfigurationId
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
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/",
			Error: true,
		},

		{
			// missing ApplicationGatewayName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for ApplicationGatewayName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/",
			Error: true,
		},

		{
			// missing PrivateLinkConfigurationName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/",
			Error: true,
		},

		{
			// missing value for PrivateLinkConfigurationName
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/privateLinkConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/00d16aa6-089c-400f-98eb-926adb838ee1/resourceGroups/cbuk-core-testoce-appgatewayshared-uksouth/providers/Microsoft.Network/applicationGateways/cbuk-core-testoce-appgatewayshared-uksouth/privateLinkConfigurations/hbtest",
			Expected: &ApplicationGatewayPrivateLinkConfigurationId{
				SubscriptionId:               "00d16aa6-089c-400f-98eb-926adb838ee1",
				ResourceGroup:                "cbuk-core-testoce-appgatewayshared-uksouth",
				ApplicationGatewayName:       "cbuk-core-testoce-appgatewayshared-uksouth",
				PrivateLinkConfigurationName: "hbtest",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00D16AA6-089C-400F-98EB-926ADB838EE1/RESOURCEGROUPS/CBUK-CORE-TESTOCE-APPGATEWAYSHARED-UKSOUTH/PROVIDERS/MICROSOFT.NETWORK/APPLICATIONGATEWAYS/CBUK-CORE-TESTOCE-APPGATEWAYSHARED-UKSOUTH/PRIVATELINKCONFIGURATIONS/HBTEST",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ApplicationGatewayPrivateLinkConfigurationID(v.Input)
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
		if actual.ApplicationGatewayName != v.Expected.ApplicationGatewayName {
			t.Fatalf("Expected %q but got %q for ApplicationGatewayName", v.Expected.ApplicationGatewayName, actual.ApplicationGatewayName)
		}
		if actual.PrivateLinkConfigurationName != v.Expected.PrivateLinkConfigurationName {
			t.Fatalf("Expected %q but got %q for PrivateLinkConfigurationName", v.Expected.PrivateLinkConfigurationName, actual.PrivateLinkConfigurationName)
		}
	}
}
