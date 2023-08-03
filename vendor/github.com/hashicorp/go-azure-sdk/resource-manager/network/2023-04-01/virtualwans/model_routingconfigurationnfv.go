package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingConfigurationNfv struct {
	AssociatedRouteTable  *RoutingConfigurationNfvSubResource `json:"associatedRouteTable,omitempty"`
	InboundRouteMap       *RoutingConfigurationNfvSubResource `json:"inboundRouteMap,omitempty"`
	OutboundRouteMap      *RoutingConfigurationNfvSubResource `json:"outboundRouteMap,omitempty"`
	PropagatedRouteTables *PropagatedRouteTableNfv            `json:"propagatedRouteTables,omitempty"`
}
