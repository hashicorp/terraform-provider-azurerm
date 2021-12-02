package parse

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = PrivateLinkConfigurationId{}

func TestPrivateLinkConfigurationIDFormatter(t *testing.T) {
	actual := NewPrivateLinkConfigurationID("12345678-1234-9876-4563-123456789012", "resGroup1", "appgw1", "plc1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/appgw1/privateLinkConfigurations/plc1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPrivateLinkConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateLinkConfigurationId
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
			// missing ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},
		{
			// missing value for ApplicationGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/",
			Error: true,
		},
		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/appgw1/",
			Error: true,
		},
		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/appgw1/privateLinkConfigurations/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/appgw1/privateLinkConfigurations/plc1",
			Expected: &PrivateLinkConfigurationId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "resGroup1",
				ApplicationGatewayName: "appgw1",
				Name:                   "plc1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/APPLICATIONGATEWAYS/APPGW1/PRIVATELINKCONFIGURATIONS/PLC1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PrivateLinkConfigurationID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
