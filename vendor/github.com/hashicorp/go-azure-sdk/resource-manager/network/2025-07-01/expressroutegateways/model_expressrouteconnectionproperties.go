package expressroutegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteConnectionProperties struct {
	AuthorizationKey           *string                      `json:"authorizationKey,omitempty"`
	EnableInternetSecurity     *bool                        `json:"enableInternetSecurity,omitempty"`
	EnablePrivateLinkFastPath  *bool                        `json:"enablePrivateLinkFastPath,omitempty"`
	ExpressRouteCircuitPeering ExpressRouteCircuitPeeringId `json:"expressRouteCircuitPeering"`
	ExpressRouteGatewayBypass  *bool                        `json:"expressRouteGatewayBypass,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	RoutingConfiguration       *RoutingConfiguration        `json:"routingConfiguration,omitempty"`
	RoutingWeight              *int64                       `json:"routingWeight,omitempty"`
}
