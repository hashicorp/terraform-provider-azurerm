package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = NetworkExpressRouteCircuitConnectionId{}

func TestNetworkExpressRouteCircuitConnectionIDFormatter(t *testing.T) {
	actual := NewExpressRouteCircuitConnectionID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "circuit1", "peering1", "connection1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNetworkExpressRouteCircuitConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NetworkExpressRouteCircuitConnectionId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// missing subscriptions
			Input: "/",
			Error: true,
		},
		{
			// missing value for subscriptions
			Input: "/subscriptions/",
			Error: true,
		},
		{
			// missing resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},
		{
			// missing value for resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},
		{
			// missing expressRouteCircuits
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/",
			Error: true,
		},
		{
			// missing value for expressRouteCircuits
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/",
			Error: true,
		},
		{
			// missing peerings
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/",
			Error: true,
		},
		{
			// missing value for peerings
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/",
			Error: true,
		},
		{
			// missing connections
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/",
			Error: true,
		},
		{
			// missing value for connections
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1",
			Expected: &NetworkExpressRouteCircuitConnectionId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				CircuitName:    "circuit1",
				PeeringName:    "peering1",
				Name:           "connection1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.NETWORK/EXPRESSROUTECIRCUITS/CIRCUIT1/PEERINGS/PEERING1/CONNECTIONS/CONNECTION1",
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
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.CircuitName != v.Expected.CircuitName {
			t.Fatalf("Expected %q but got %q for CircuitName", v.Expected.CircuitName, actual.CircuitName)
		}

		if actual.PeeringName != v.Expected.PeeringName {
			t.Fatalf("Expected %q but got %q for PeeringName", v.Expected.PeeringName, actual.PeeringName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
