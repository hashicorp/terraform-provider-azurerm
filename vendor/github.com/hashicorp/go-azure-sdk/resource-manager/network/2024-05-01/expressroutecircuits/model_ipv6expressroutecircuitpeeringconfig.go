package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPv6ExpressRouteCircuitPeeringConfig struct {
	MicrosoftPeeringConfig     *ExpressRouteCircuitPeeringConfig `json:"microsoftPeeringConfig,omitempty"`
	PrimaryPeerAddressPrefix   *string                           `json:"primaryPeerAddressPrefix,omitempty"`
	RouteFilter                *SubResource                      `json:"routeFilter,omitempty"`
	SecondaryPeerAddressPrefix *string                           `json:"secondaryPeerAddressPrefix,omitempty"`
	State                      *ExpressRouteCircuitPeeringState  `json:"state,omitempty"`
}
