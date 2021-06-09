package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = PrivateDnsZoneConfigId{}

func TestPrivateDnsZoneConfigIDFormatter(t *testing.T) {
	actual := NewPrivateDnsZoneConfigID("12345678-1234-9876-4563-123456789012", "resGroup1", "endpoint1", "privateDnsZoneGroup1", "privateDnsZoneConfig1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPrivateDnsZoneConfigID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateDnsZoneConfigId
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
			// missing PrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for PrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/",
			Error: true,
		},

		{
			// missing PrivateDnsZoneGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/",
			Error: true,
		},

		{
			// missing value for PrivateDnsZoneGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1",
			Expected: &PrivateDnsZoneConfigId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "resGroup1",
				PrivateEndpointName:     "endpoint1",
				PrivateDnsZoneGroupName: "privateDnsZoneGroup1",
				Name:                    "privateDnsZoneConfig1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/PRIVATEENDPOINTS/ENDPOINT1/PRIVATEDNSZONEGROUPS/PRIVATEDNSZONEGROUP1/PRIVATEDNSZONECONFIGS/PRIVATEDNSZONECONFIG1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PrivateDnsZoneConfigID(v.Input)
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
		if actual.PrivateEndpointName != v.Expected.PrivateEndpointName {
			t.Fatalf("Expected %q but got %q for PrivateEndpointName", v.Expected.PrivateEndpointName, actual.PrivateEndpointName)
		}
		if actual.PrivateDnsZoneGroupName != v.Expected.PrivateDnsZoneGroupName {
			t.Fatalf("Expected %q but got %q for PrivateDnsZoneGroupName", v.Expected.PrivateDnsZoneGroupName, actual.PrivateDnsZoneGroupName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
