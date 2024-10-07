package expressrouteproviderports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteProviderPortProperties struct {
	OverprovisionFactor      *int64  `json:"overprovisionFactor,omitempty"`
	PeeringLocation          *string `json:"peeringLocation,omitempty"`
	PortBandwidthInMbps      *int64  `json:"portBandwidthInMbps,omitempty"`
	PortPairDescriptor       *string `json:"portPairDescriptor,omitempty"`
	PrimaryAzurePort         *string `json:"primaryAzurePort,omitempty"`
	RemainingBandwidthInMbps *int64  `json:"remainingBandwidthInMbps,omitempty"`
	SecondaryAzurePort       *string `json:"secondaryAzurePort,omitempty"`
	UsedBandwidthInMbps      *int64  `json:"usedBandwidthInMbps,omitempty"`
}
