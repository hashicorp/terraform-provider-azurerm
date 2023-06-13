package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterLoadBalancerProfile struct {
	AllocatedOutboundPorts              *int64                                               `json:"allocatedOutboundPorts,omitempty"`
	BackendPoolType                     *BackendPoolType                                     `json:"backendPoolType,omitempty"`
	EffectiveOutboundIPs                *[]ResourceReference                                 `json:"effectiveOutboundIPs,omitempty"`
	EnableMultipleStandardLoadBalancers *bool                                                `json:"enableMultipleStandardLoadBalancers,omitempty"`
	IdleTimeoutInMinutes                *int64                                               `json:"idleTimeoutInMinutes,omitempty"`
	ManagedOutboundIPs                  *ManagedClusterLoadBalancerProfileManagedOutboundIPs `json:"managedOutboundIPs,omitempty"`
	OutboundIPPrefixes                  *ManagedClusterLoadBalancerProfileOutboundIPPrefixes `json:"outboundIPPrefixes,omitempty"`
	OutboundIPs                         *ManagedClusterLoadBalancerProfileOutboundIPs        `json:"outboundIPs,omitempty"`
}
