package virtualnetworktaps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceIPConfigurationPropertiesFormat struct {
	ApplicationGatewayBackendAddressPools *[]ApplicationGatewayBackendAddressPool                         `json:"applicationGatewayBackendAddressPools,omitempty"`
	ApplicationSecurityGroups             *[]ApplicationSecurityGroup                                     `json:"applicationSecurityGroups,omitempty"`
	GatewayLoadBalancer                   *SubResource                                                    `json:"gatewayLoadBalancer,omitempty"`
	LoadBalancerBackendAddressPools       *[]BackendAddressPool                                           `json:"loadBalancerBackendAddressPools,omitempty"`
	LoadBalancerInboundNatRules           *[]InboundNatRule                                               `json:"loadBalancerInboundNatRules,omitempty"`
	Primary                               *bool                                                           `json:"primary,omitempty"`
	PrivateIPAddress                      *string                                                         `json:"privateIPAddress,omitempty"`
	PrivateIPAddressVersion               *IPVersion                                                      `json:"privateIPAddressVersion,omitempty"`
	PrivateIPAllocationMethod             *IPAllocationMethod                                             `json:"privateIPAllocationMethod,omitempty"`
	PrivateLinkConnectionProperties       *NetworkInterfaceIPConfigurationPrivateLinkConnectionProperties `json:"privateLinkConnectionProperties,omitempty"`
	ProvisioningState                     *ProvisioningState                                              `json:"provisioningState,omitempty"`
	PublicIPAddress                       *PublicIPAddress                                                `json:"publicIPAddress,omitempty"`
	Subnet                                *Subnet                                                         `json:"subnet,omitempty"`
	VirtualNetworkTaps                    *[]VirtualNetworkTap                                            `json:"virtualNetworkTaps,omitempty"`
}
