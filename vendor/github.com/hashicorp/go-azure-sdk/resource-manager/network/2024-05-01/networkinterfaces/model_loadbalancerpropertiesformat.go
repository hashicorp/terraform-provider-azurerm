package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerPropertiesFormat struct {
	BackendAddressPools      *[]BackendAddressPool      `json:"backendAddressPools,omitempty"`
	FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
	InboundNatPools          *[]InboundNatPool          `json:"inboundNatPools,omitempty"`
	InboundNatRules          *[]InboundNatRule          `json:"inboundNatRules,omitempty"`
	LoadBalancingRules       *[]LoadBalancingRule       `json:"loadBalancingRules,omitempty"`
	OutboundRules            *[]OutboundRule            `json:"outboundRules,omitempty"`
	Probes                   *[]Probe                   `json:"probes,omitempty"`
	ProvisioningState        *ProvisioningState         `json:"provisioningState,omitempty"`
	ResourceGuid             *string                    `json:"resourceGuid,omitempty"`
}
