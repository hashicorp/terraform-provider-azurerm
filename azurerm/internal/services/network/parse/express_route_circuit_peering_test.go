package parse

import (
	"testing"
)

func TestExpressRouteCircuitPeeringID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ExpressRouteCircuitPeeringId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Peering Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings",
			Expected: nil,
		},
		{
			Name:  "network ExpressRouteCircuitPeering ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1",
			Expected: &ExpressRouteCircuitPeeringId{
				ResourceGroup: "resourceGroup1",
				CircuitName:   "circuit1",
				Name:          "peering1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/Peerings/peering1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := ExpressRouteCircuitPeeringID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.CircuitName != v.Expected.CircuitName {
			t.Fatalf("Expected %q but got %q for CircuitName", v.Expected.CircuitName, actual.CircuitName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
