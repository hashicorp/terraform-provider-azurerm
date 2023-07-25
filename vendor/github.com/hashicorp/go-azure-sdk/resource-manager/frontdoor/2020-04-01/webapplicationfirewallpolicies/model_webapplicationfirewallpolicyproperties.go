package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallPolicyProperties struct {
	CustomRules           *CustomRuleList         `json:"customRules,omitempty"`
	FrontendEndpointLinks *[]FrontendEndpointLink `json:"frontendEndpointLinks,omitempty"`
	ManagedRules          *ManagedRuleSetList     `json:"managedRules,omitempty"`
	PolicySettings        *PolicySettings         `json:"policySettings,omitempty"`
	ProvisioningState     *string                 `json:"provisioningState,omitempty"`
	ResourceState         *PolicyResourceState    `json:"resourceState,omitempty"`
	RoutingRuleLinks      *[]RoutingRuleLink      `json:"routingRuleLinks,omitempty"`
}
