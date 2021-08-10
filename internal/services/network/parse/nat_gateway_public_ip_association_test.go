package parse

import (
	"testing"
)

func TestNatGatewayPublicIPAddressAssociationID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *NatGatewayPublicIPAddressAssociationId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "One Segment",
			Input: "hello",
			Error: true,
		},
		{
			Name:  "Two Segments Invalid ID's",
			Input: "hello|world",
			Error: true,
		},
		{
			Name:  "Missing Nat Gateway Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways",
			Error: true,
		},
		{
			Name:  "Nat Gateway ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1",
			Error: true,
		},
		{
			Name:  "Public IP Address ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1",
			Error: true,
		},
		{
			Name:  "Nat Gateway / Public IP Association ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1",
			Error: false,
			Expect: &NatGatewayPublicIPAddressAssociationId{
				NatGateway: NatGatewayId{
					Name:          "gateway1",
					ResourceGroup: "group1",
				},
				PublicIPAddressID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := NatGatewayPublicIPAddressAssociationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.NatGateway.Name != v.Expect.NatGateway.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.NatGateway.Name, actual.NatGateway.Name)
		}

		if actual.NatGateway.ResourceGroup != v.Expect.NatGateway.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.NatGateway.ResourceGroup, actual.NatGateway.ResourceGroup)
		}
	}
}
