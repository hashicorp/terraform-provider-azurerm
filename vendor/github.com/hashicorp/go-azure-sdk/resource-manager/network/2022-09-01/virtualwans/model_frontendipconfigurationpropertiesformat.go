package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendIPConfigurationPropertiesFormat struct {
	GatewayLoadBalancer       *SubResource        `json:"gatewayLoadBalancer,omitempty"`
	InboundNatPools           *[]SubResource      `json:"inboundNatPools,omitempty"`
	InboundNatRules           *[]SubResource      `json:"inboundNatRules,omitempty"`
	LoadBalancingRules        *[]SubResource      `json:"loadBalancingRules,omitempty"`
	OutboundRules             *[]SubResource      `json:"outboundRules,omitempty"`
	PrivateIPAddress          *string             `json:"privateIPAddress,omitempty"`
	PrivateIPAddressVersion   *IPVersion          `json:"privateIPAddressVersion,omitempty"`
	PrivateIPAllocationMethod *IPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	ProvisioningState         *ProvisioningState  `json:"provisioningState,omitempty"`
	PublicIPAddress           *PublicIPAddress    `json:"publicIPAddress,omitempty"`
	PublicIPPrefix            *SubResource        `json:"publicIPPrefix,omitempty"`
	Subnet                    *Subnet             `json:"subnet,omitempty"`
}
