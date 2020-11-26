package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VPNGatewayId{}

func TestVPNGatewayIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewVPNGatewayID("group1", "gateway1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseVPNGateway(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VPNGatewayId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No VPN Gateways Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No VPN Gateways Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnGateways/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/vpnGateways/example",
			Expected: &VPNGatewayId{
				Name:          "example",
				ResourceGroup: "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VPNGatewayID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

var _ resourceid.Formatter = VPNGatewayConnectionId{}

func TestVPNGatewayConnectionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewVPNGatewayConnectionID("group1", "gateway1", "conn1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/conn1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVPNGatewayConnectionID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *VPNGatewayConnectionId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Missing leading slash",
			Input: "subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1",
			Error: true,
		},
		{
			Name:  "Malformed segments",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/foo/bar",
			Error: true,
		},
		{
			Name:  "No vpn gateway segment",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/providers/Microsoft.Network",
			Error: true,
		},
		{
			Name:  "No vpn connection segment",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1",
			Error: true,
		},
		{
			Name:  "Correct",
			Input: "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/group1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/conn1",
			Expect: &VPNGatewayConnectionId{
				ResourceGroup: "group1",
				Gateway:       "gateway1",
				Name:          "conn1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VPNGatewayConnectionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get")
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Gateway != v.Expect.Gateway {
			t.Fatalf("Expected %q but got %q for Gateway", v.Expect.Gateway, actual.Gateway)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
