package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServiceNetworkProfile struct {
	DnsServiceIP        *string                                        `json:"dnsServiceIP,omitempty"`
	DockerBridgeCidr    *string                                        `json:"dockerBridgeCidr,omitempty"`
	EbpfDataplane       *EbpfDataplane                                 `json:"ebpfDataplane,omitempty"`
	IPFamilies          *[]IPFamily                                    `json:"ipFamilies,omitempty"`
	KubeProxyConfig     *ContainerServiceNetworkProfileKubeProxyConfig `json:"kubeProxyConfig,omitempty"`
	LoadBalancerProfile *ManagedClusterLoadBalancerProfile             `json:"loadBalancerProfile,omitempty"`
	LoadBalancerSku     *LoadBalancerSku                               `json:"loadBalancerSku,omitempty"`
	NatGatewayProfile   *ManagedClusterNATGatewayProfile               `json:"natGatewayProfile,omitempty"`
	NetworkMode         *NetworkMode                                   `json:"networkMode,omitempty"`
	NetworkPlugin       *NetworkPlugin                                 `json:"networkPlugin,omitempty"`
	NetworkPluginMode   *NetworkPluginMode                             `json:"networkPluginMode,omitempty"`
	NetworkPolicy       *NetworkPolicy                                 `json:"networkPolicy,omitempty"`
	OutboundType        *OutboundType                                  `json:"outboundType,omitempty"`
	PodCidr             *string                                        `json:"podCidr,omitempty"`
	PodCidrs            *[]string                                      `json:"podCidrs,omitempty"`
	ServiceCidr         *string                                        `json:"serviceCidr,omitempty"`
	ServiceCidrs        *[]string                                      `json:"serviceCidrs,omitempty"`
}
