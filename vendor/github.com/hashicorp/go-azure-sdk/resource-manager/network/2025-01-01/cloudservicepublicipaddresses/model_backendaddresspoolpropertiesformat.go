package cloudservicepublicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendAddressPoolPropertiesFormat struct {
	BackendIPConfigurations      *[]NetworkInterfaceIPConfiguration    `json:"backendIPConfigurations,omitempty"`
	DrainPeriodInSeconds         *int64                                `json:"drainPeriodInSeconds,omitempty"`
	InboundNatRules              *[]SubResource                        `json:"inboundNatRules,omitempty"`
	LoadBalancerBackendAddresses *[]LoadBalancerBackendAddress         `json:"loadBalancerBackendAddresses,omitempty"`
	LoadBalancingRules           *[]SubResource                        `json:"loadBalancingRules,omitempty"`
	Location                     *string                               `json:"location,omitempty"`
	OutboundRule                 *SubResource                          `json:"outboundRule,omitempty"`
	OutboundRules                *[]SubResource                        `json:"outboundRules,omitempty"`
	ProvisioningState            *ProvisioningState                    `json:"provisioningState,omitempty"`
	SyncMode                     *SyncMode                             `json:"syncMode,omitempty"`
	TunnelInterfaces             *[]GatewayLoadBalancerTunnelInterface `json:"tunnelInterfaces,omitempty"`
	VirtualNetwork               *SubResource                          `json:"virtualNetwork,omitempty"`
}
