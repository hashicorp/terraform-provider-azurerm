package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ExpressRouteCircuitConnectionId{}

func TestExpressRouteCircuitConnectionIDFormatter(t *testing.T) {
	actual := NewExpressRouteCircuitConnectionID("12345678-1234-9876-4563-123456789012", "resGroup1", "circuit1", "peering1", "connection1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestExpressRouteCircuitConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ExpressRouteCircuitConnectionId
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
			// missing ExpressRouteCircuitName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for ExpressRouteCircuitName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/",
			Error: true,
		},

		{
			// missing PeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/",
			Error: true,
		},

		{
			// missing value for PeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/",
			Error: true,
		},

		{
			// missing ConnectionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/",
			Error: true,
		},

		{
			// missing value for ConnectionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1",
			Expected: &ExpressRouteCircuitConnectionId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "resGroup1",
				ExpressRouteCircuitName: "circuit1",
				PeeringName:             "peering1",
				ConnectionName:          "connection1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/EXPRESSROUTECIRCUITS/CIRCUIT1/PEERINGS/PEERING1/CONNECTIONS/CONNECTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ExpressRouteCircuitConnectionID(v.Input)
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
		if actual.ExpressRouteCircuitName != v.Expected.ExpressRouteCircuitName {
			t.Fatalf("Expected %q but got %q for ExpressRouteCircuitName", v.Expected.ExpressRouteCircuitName, actual.ExpressRouteCircuitName)
		}
		if actual.PeeringName != v.Expected.PeeringName {
			t.Fatalf("Expected %q but got %q for PeeringName", v.Expected.PeeringName, actual.PeeringName)
		}
		if actual.ConnectionName != v.Expected.ConnectionName {
			t.Fatalf("Expected %q but got %q for ConnectionName", v.Expected.ConnectionName, actual.ConnectionName)
		}
	}
}
