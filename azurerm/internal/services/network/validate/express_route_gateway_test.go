package validate

import "testing"

func TestExpressRouteGatewayID(t *testing.T) {
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
			Name:  "No expressRouteGateways Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			Name:  "No expressRouteGateways Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteGateways/",
			Valid: false,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/expressRouteGateways/expressRouteGateway1",
			Valid: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		_, errors := ExpressRouteGatewayID(v.Input, "express_route_gateway_id")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
