package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfiguration struct {
	ApplicationGatewayBackendAddressPools *[]SubResource                               `json:"applicationGatewayBackendAddressPools,omitempty"`
	LoadBalancerBackendAddressPools       *[]SubResource                               `json:"loadBalancerBackendAddressPools,omitempty"`
	LoadBalancerInboundNatPools           *[]SubResource                               `json:"loadBalancerInboundNatPools,omitempty"`
	Name                                  string                                       `json:"name"`
	PrivateIPAddressVersion               *PrivateIPAddressVersion                     `json:"privateIPAddressVersion,omitempty"`
	PublicIPAddressConfiguration          *IPConfigurationPublicIPAddressConfiguration `json:"publicIPAddressConfiguration,omitempty"`
	Subnet                                *SubResource                                 `json:"subnet,omitempty"`
}
