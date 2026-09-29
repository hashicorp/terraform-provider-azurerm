// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworkgatewayconnections"
)

func TestConnectionTypeUsesSharedKey(t *testing.T) {
	testData := []struct {
		name           string
		connectionType virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionType
		expected       bool
	}{
		{
			name:           "IPsec uses a Shared Key",
			connectionType: virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec,
			expected:       true,
		},
		{
			name:           "Vnet2Vnet uses a Shared Key",
			connectionType: virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet,
			expected:       true,
		},
		{
			name:           "ExpressRoute uses an Authorization Key, not a Shared Key",
			connectionType: virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute,
			expected:       false,
		},
		{
			name:           "an empty connection type is not assumed to be ExpressRoute",
			connectionType: "",
			expected:       true,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			if actual := connectionTypeUsesSharedKey(v.connectionType); actual != v.expected {
				t.Fatalf("expected %t for connection type %q but got %t", v.expected, string(v.connectionType), actual)
			}
		})
	}
}
