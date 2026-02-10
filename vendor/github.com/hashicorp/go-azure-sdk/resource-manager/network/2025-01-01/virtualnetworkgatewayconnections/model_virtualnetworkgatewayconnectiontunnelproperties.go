package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayConnectionTunnelProperties struct {
	BgpPeeringAddress *string `json:"bgpPeeringAddress,omitempty"`
	TunnelIPAddress   *string `json:"tunnelIpAddress,omitempty"`
}
