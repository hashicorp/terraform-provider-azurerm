package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = ExpressRouteGatewayId{}

func TestExpressRouteGatewayIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	id := NewExpressRouteGatewayID("resourceGroup1", "expressRouteGateway1")
	actual := id.ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteGateways/expressRouteGateway1"

	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestExpressRouteGatewayID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ExpressRouteGatewayId
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
			Name:     "Missing ExpressRouteGateway Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteGateways",
			Expected: nil,
		},
		{
			Name:  "network ExpressRouteGateway ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/expressRouteGateways/expressRouteGateway1",
			Expected: &ExpressRouteGatewayId{
				ResourceGroup: "resourceGroup1",
				Name:          "expressRouteGateway1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/ExpressRouteGateways/expressRouteGateway1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := ExpressRouteGatewayID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
