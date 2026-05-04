package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayURLPathMapPropertiesFormat struct {
	DefaultBackendAddressPool     *SubResource                  `json:"defaultBackendAddressPool,omitempty"`
	DefaultBackendHTTPSettings    *SubResource                  `json:"defaultBackendHttpSettings,omitempty"`
	DefaultLoadDistributionPolicy *SubResource                  `json:"defaultLoadDistributionPolicy,omitempty"`
	DefaultRedirectConfiguration  *SubResource                  `json:"defaultRedirectConfiguration,omitempty"`
	DefaultRewriteRuleSet         *SubResource                  `json:"defaultRewriteRuleSet,omitempty"`
	PathRules                     *[]ApplicationGatewayPathRule `json:"pathRules,omitempty"`
	ProvisioningState             *ProvisioningState            `json:"provisioningState,omitempty"`
}
