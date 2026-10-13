package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingConfiguration struct {
	AssociatedRouteTable  *SubResource          `json:"associatedRouteTable,omitempty"`
	InboundRouteMap       *SubResource          `json:"inboundRouteMap,omitempty"`
	OutboundRouteMap      *SubResource          `json:"outboundRouteMap,omitempty"`
	PropagatedRouteTables *PropagatedRouteTable `json:"propagatedRouteTables,omitempty"`
	VnetRoutes            *VnetRoute            `json:"vnetRoutes,omitempty"`
}
