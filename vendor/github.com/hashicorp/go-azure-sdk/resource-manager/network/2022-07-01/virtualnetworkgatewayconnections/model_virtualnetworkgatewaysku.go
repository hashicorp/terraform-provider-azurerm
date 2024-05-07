package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewaySku struct {
	Capacity *int64                        `json:"capacity,omitempty"`
	Name     *VirtualNetworkGatewaySkuName `json:"name,omitempty"`
	Tier     *VirtualNetworkGatewaySkuTier `json:"tier,omitempty"`
}
