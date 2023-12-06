package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitPeeringConfig struct {
	AdvertisedCommunities         *[]string                                              `json:"advertisedCommunities,omitempty"`
	AdvertisedPublicPrefixes      *[]string                                              `json:"advertisedPublicPrefixes,omitempty"`
	AdvertisedPublicPrefixesState *ExpressRouteCircuitPeeringAdvertisedPublicPrefixState `json:"advertisedPublicPrefixesState,omitempty"`
	CustomerASN                   *int64                                                 `json:"customerASN,omitempty"`
	LegacyMode                    *int64                                                 `json:"legacyMode,omitempty"`
	RoutingRegistryName           *string                                                `json:"routingRegistryName,omitempty"`
}
