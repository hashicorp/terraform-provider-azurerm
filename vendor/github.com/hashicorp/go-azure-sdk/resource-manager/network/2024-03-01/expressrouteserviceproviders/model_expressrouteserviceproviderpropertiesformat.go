package expressrouteserviceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteServiceProviderPropertiesFormat struct {
	BandwidthsOffered *[]ExpressRouteServiceProviderBandwidthsOffered `json:"bandwidthsOffered,omitempty"`
	PeeringLocations  *[]string                                       `json:"peeringLocations,omitempty"`
	ProvisioningState *ProvisioningState                              `json:"provisioningState,omitempty"`
}
