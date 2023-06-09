package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancingRulePropertiesFormat struct {
	BackendAddressPool      *SubResource       `json:"backendAddressPool,omitempty"`
	BackendAddressPools     *[]SubResource     `json:"backendAddressPools,omitempty"`
	BackendPort             *int64             `json:"backendPort,omitempty"`
	DisableOutboundSnat     *bool              `json:"disableOutboundSnat,omitempty"`
	EnableFloatingIP        *bool              `json:"enableFloatingIP,omitempty"`
	EnableTcpReset          *bool              `json:"enableTcpReset,omitempty"`
	FrontendIPConfiguration *SubResource       `json:"frontendIPConfiguration,omitempty"`
	FrontendPort            int64              `json:"frontendPort"`
	IdleTimeoutInMinutes    *int64             `json:"idleTimeoutInMinutes,omitempty"`
	LoadDistribution        *LoadDistribution  `json:"loadDistribution,omitempty"`
	Probe                   *SubResource       `json:"probe,omitempty"`
	Protocol                TransportProtocol  `json:"protocol"`
	ProvisioningState       *ProvisioningState `json:"provisioningState,omitempty"`
}
