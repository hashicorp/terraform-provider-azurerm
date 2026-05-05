package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundNatPoolPropertiesFormat struct {
	BackendPort             int64              `json:"backendPort"`
	EnableFloatingIP        *bool              `json:"enableFloatingIP,omitempty"`
	EnableTcpReset          *bool              `json:"enableTcpReset,omitempty"`
	FrontendIPConfiguration *SubResource       `json:"frontendIPConfiguration,omitempty"`
	FrontendPortRangeEnd    int64              `json:"frontendPortRangeEnd"`
	FrontendPortRangeStart  int64              `json:"frontendPortRangeStart"`
	IdleTimeoutInMinutes    *int64             `json:"idleTimeoutInMinutes,omitempty"`
	Protocol                TransportProtocol  `json:"protocol"`
	ProvisioningState       *ProvisioningState `json:"provisioningState,omitempty"`
}
