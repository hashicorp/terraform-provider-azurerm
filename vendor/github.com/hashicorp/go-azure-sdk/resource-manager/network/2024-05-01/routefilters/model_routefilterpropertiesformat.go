package routefilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteFilterPropertiesFormat struct {
	IPv6Peerings      *[]ExpressRouteCircuitPeering `json:"ipv6Peerings,omitempty"`
	Peerings          *[]ExpressRouteCircuitPeering `json:"peerings,omitempty"`
	ProvisioningState *ProvisioningState            `json:"provisioningState,omitempty"`
	Rules             *[]RouteFilterRule            `json:"rules,omitempty"`
}
