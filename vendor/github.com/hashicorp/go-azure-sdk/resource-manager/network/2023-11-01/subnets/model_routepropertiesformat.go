package subnets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutePropertiesFormat struct {
	AddressPrefix     *string            `json:"addressPrefix,omitempty"`
	HasBgpOverride    *bool              `json:"hasBgpOverride,omitempty"`
	NextHopIPAddress  *string            `json:"nextHopIpAddress,omitempty"`
	NextHopType       RouteNextHopType   `json:"nextHopType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
