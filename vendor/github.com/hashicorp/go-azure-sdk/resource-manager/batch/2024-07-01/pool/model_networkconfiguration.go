package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfiguration struct {
	DynamicVnetAssignmentScope   *DynamicVNetAssignmentScope   `json:"dynamicVnetAssignmentScope,omitempty"`
	EnableAcceleratedNetworking  *bool                         `json:"enableAcceleratedNetworking,omitempty"`
	EndpointConfiguration        *PoolEndpointConfiguration    `json:"endpointConfiguration,omitempty"`
	PublicIPAddressConfiguration *PublicIPAddressConfiguration `json:"publicIPAddressConfiguration,omitempty"`
	SubnetId                     *string                       `json:"subnetId,omitempty"`
}
