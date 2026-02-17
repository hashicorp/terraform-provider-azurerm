// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"reflect"
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

func TestFlattenVirtualNetworkGatewayIPConfigurations(t *testing.T) {
	subnetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/GatewaySubnet"
	publicIPID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/pip"

	tests := []struct {
		name  string
		input *[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration
		want  []map[string]interface{}
	}{
		{
			name:  "nil input returns empty output",
			input: nil,
			want:  []map[string]interface{}{},
		},
		{
			name:  "empty input returns empty output",
			input: &[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{},
			want:  []map[string]interface{}{},
		},
		{
			name: "includes public ip address id",
			input: &[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
				{
					Name: pointer.To("vnetGatewayConfig"),
					Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethodDynamic),
						Subnet: &virtualnetworkgateways.SubResource{
							Id: pointer.To(subnetID),
						},
						PublicIPAddress: &virtualnetworkgateways.SubResource{
							Id: pointer.To(publicIPID),
						},
					},
				},
			},
			want: []map[string]interface{}{
				{
					"name":                          "vnetGatewayConfig",
					"private_ip_address_allocation": "Dynamic",
					"subnet_id":                     subnetID,
					"public_ip_address_id":          publicIPID,
				},
			},
		},
		{
			name: "omits optional fields when ids and name are absent",
			input: &[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
				{
					Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethodStatic),
					},
				},
			},
			want: []map[string]interface{}{
				{
					"private_ip_address_allocation": "Static",
				},
			},
		},
		{
			name: "preserves order for multiple configurations",
			input: &[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
				{
					Name: pointer.To("first"),
					Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethodDynamic),
					},
				},
				{
					Name: pointer.To("second"),
					Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethodStatic),
					},
				},
			},
			want: []map[string]interface{}{
				{
					"name":                          "first",
					"private_ip_address_allocation": "Dynamic",
				},
				{
					"name":                          "second",
					"private_ip_address_allocation": "Static",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotRaw := flattenVirtualNetworkGatewayIPConfigurations(tt.input)
			got := make([]map[string]interface{}, 0, len(gotRaw))

			for i, item := range gotRaw {
				cfg, ok := item.(map[string]interface{})
				if !ok {
					t.Fatalf("expected flattened item at index %d to be a map[string]interface{}", i)
				}

				got = append(got, cfg)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected flattened output: got %#v, want %#v", got, tt.want)
			}
		})
	}
}
