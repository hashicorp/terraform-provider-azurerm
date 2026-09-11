package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineNetworkInterfaceIPConfigurationProperties struct {
	ApplicationGatewayBackendAddressPools *[]SubResource                              `json:"applicationGatewayBackendAddressPools,omitempty"`
	ApplicationSecurityGroups             *[]SubResource                              `json:"applicationSecurityGroups,omitempty"`
	LoadBalancerBackendAddressPools       *[]SubResource                              `json:"loadBalancerBackendAddressPools,omitempty"`
	Primary                               *bool                                       `json:"primary,omitempty"`
	PrivateIPAddressVersion               *IPVersions                                 `json:"privateIPAddressVersion,omitempty"`
	PublicIPAddressConfiguration          *VirtualMachinePublicIPAddressConfiguration `json:"publicIPAddressConfiguration,omitempty"`
	Subnet                                *SubResource                                `json:"subnet,omitempty"`
}
