package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayEffectiveRoute struct {
	AddressPrefixes  *[]string                         `json:"addressPrefixes,omitempty"`
	LocalAddress     *string                           `json:"localAddress,omitempty"`
	NextHopIPAddress *string                           `json:"nextHopIpAddress,omitempty"`
	NextHopType      *GatewayEffectiveRouteNextHopType `json:"nextHopType,omitempty"`
}
