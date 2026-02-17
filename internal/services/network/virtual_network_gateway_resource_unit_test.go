// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworkgateways"
)

func TestVirtualNetworkGatewayResource_NoCustomizeDiff(t *testing.T) {
	r := resourceVirtualNetworkGateway()

	if r.CustomizeDiff != nil {
		t.Fatalf("expected CustomizeDiff to be nil")
	}
}

func TestFlattenVirtualNetworkGatewayIPConfigurations_IncludesPublicIPAddress(t *testing.T) {
	input := &[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
		{
			Name: pointer.To("vnetGatewayConfig"),
			Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethodDynamic),
				Subnet: &virtualnetworkgateways.SubResource{
					Id: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/GatewaySubnet"),
				},
				PublicIPAddress: &virtualnetworkgateways.SubResource{
					Id: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/pip"),
				},
			},
		},
	}

	output := flattenVirtualNetworkGatewayIPConfigurations(input)
	if len(output) != 1 {
		t.Fatalf("expected one flattened ip_configuration but got %d", len(output))
	}

	config, ok := output[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected flattened item to be a map[string]interface{}")
	}

	publicIPID, ok := config["public_ip_address_id"].(string)
	if !ok {
		t.Fatalf("expected public_ip_address_id to be present in flattened output")
	}

	expectedPublicIPID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/pip"
	if publicIPID != expectedPublicIPID {
		t.Fatalf("expected public_ip_address_id %q but got %q", expectedPublicIPID, publicIPID)
	}
}

