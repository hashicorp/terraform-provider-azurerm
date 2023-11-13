package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallPolicyPropertiesFormat struct {
	ApplicationGateways *[]ApplicationGateway                      `json:"applicationGateways,omitempty"`
	CustomRules         *[]WebApplicationFirewallCustomRule        `json:"customRules,omitempty"`
	HTTPListeners       *[]SubResource                             `json:"httpListeners,omitempty"`
	ManagedRules        ManagedRulesDefinition                     `json:"managedRules"`
	PathBasedRules      *[]SubResource                             `json:"pathBasedRules,omitempty"`
	PolicySettings      *PolicySettings                            `json:"policySettings,omitempty"`
	ProvisioningState   *ProvisioningState                         `json:"provisioningState,omitempty"`
	ResourceState       *WebApplicationFirewallPolicyResourceState `json:"resourceState,omitempty"`
}
