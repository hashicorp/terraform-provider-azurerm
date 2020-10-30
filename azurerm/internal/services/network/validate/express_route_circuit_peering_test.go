package validate

import "testing"

func TestExpressRouteCircuitPeeringID(t *testing.T) {
	testData := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Valid: false,
		},
		{
			Name:  "No expressRouteCircuits Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Name:  "No expressRouteCircuits Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteCircuits/",
			Valid: false,
		},
		{
			Name:  "No peerings Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteCircuits/circuit1",
			Valid: false,
		},
		{
			Name:  "No peerings Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteCircuits/circuit1/peerings/",
			Valid: false,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteCircuits/circuit1/peerings/peering1",
			Valid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		_, errors := ExpressRouteCircuitPeeringID(v.Input, "express_route_circuit_peering_id")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
