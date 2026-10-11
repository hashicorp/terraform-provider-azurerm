package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayPathRulePropertiesFormat struct {
	BackendAddressPool     *SubResource       `json:"backendAddressPool,omitempty"`
	BackendHTTPSettings    *SubResource       `json:"backendHttpSettings,omitempty"`
	FirewallPolicy         *SubResource       `json:"firewallPolicy,omitempty"`
	LoadDistributionPolicy *SubResource       `json:"loadDistributionPolicy,omitempty"`
	Paths                  *[]string          `json:"paths,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	RedirectConfiguration  *SubResource       `json:"redirectConfiguration,omitempty"`
	RewriteRuleSet         *SubResource       `json:"rewriteRuleSet,omitempty"`
}
