// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworkgateways"
)

func TestFlattenVirtualNetworkGatewayIPConfigurations_PublicIPAddressID(t *testing.T) {
	validPublicIPAddressID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/test-pip"
	privateIPAllocationMethod := virtualnetworkgateways.IPAllocationMethodDynamic

	t.Run("normalizes public ip id via common id parser", func(t *testing.T) {
		input := []virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
			{
				Name: pointer.To("cfg-1"),
				Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
					PrivateIPAllocationMethod: pointer.To(privateIPAllocationMethod),
					PublicIPAddress: &virtualnetworkgateways.SubResource{
						Id: pointer.To(validPublicIPAddressID),
					},
				},
			},
		}

		actual, err := flattenVirtualNetworkGatewayIPConfigurations(&input)
		if err != nil {
			t.Fatalf("expected no error but got %s", err)
		}

		if len(actual) != 1 {
			t.Fatalf("expected one flattened config but got %d", len(actual))
		}

		config, ok := actual[0].(map[string]interface{})
		if !ok {
			t.Fatalf("expected map[string]interface{} but got %T", actual[0])
		}

		gotPublicIPAddressID, ok := config["public_ip_address_id"].(string)
		if !ok {
			t.Fatalf("expected public_ip_address_id to be a string but got %T", config["public_ip_address_id"])
		}
		if gotPublicIPAddressID != validPublicIPAddressID {
			t.Fatalf("expected public_ip_address_id to be %q but got %q", validPublicIPAddressID, gotPublicIPAddressID)
		}
	})

	t.Run("returns an error when public ip id is invalid", func(t *testing.T) {
		input := []virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
			{
				Name: pointer.To("cfg-1"),
				Properties: &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
					PrivateIPAllocationMethod: pointer.To(privateIPAllocationMethod),
					PublicIPAddress: &virtualnetworkgateways.SubResource{
						Id: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses"),
					},
				},
			},
		}

		_, err := flattenVirtualNetworkGatewayIPConfigurations(&input)
		if err == nil {
			t.Fatal("expected an error but got nil")
		}
	})
}
