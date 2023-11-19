package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerBackendAddressPropertiesFormat struct {
	AdminState                          *LoadBalancerBackendAddressAdminState `json:"adminState,omitempty"`
	IPAddress                           *string                               `json:"ipAddress,omitempty"`
	InboundNatRulesPortMapping          *[]NatRulePortMapping                 `json:"inboundNatRulesPortMapping,omitempty"`
	LoadBalancerFrontendIPConfiguration *SubResource                          `json:"loadBalancerFrontendIPConfiguration,omitempty"`
	NetworkInterfaceIPConfiguration     *SubResource                          `json:"networkInterfaceIPConfiguration,omitempty"`
	Subnet                              *SubResource                          `json:"subnet,omitempty"`
	VirtualNetwork                      *SubResource                          `json:"virtualNetwork,omitempty"`
}
