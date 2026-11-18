package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAdvancedRoutingMapPropertiesFormat struct {
	AdvancedRoutingRules         []ApplicationGatewayAdvancedRoutingRule `json:"advancedRoutingRules"`
	DefaultAuthConfigs           *[]ApplicationGatewayAuthConfig         `json:"defaultAuthConfigs,omitempty"`
	DefaultBackendAddressPool    *SubResource                            `json:"defaultBackendAddressPool,omitempty"`
	DefaultBackendHTTPSettings   *SubResource                            `json:"defaultBackendHttpSettings,omitempty"`
	DefaultRedirectConfiguration *SubResource                            `json:"defaultRedirectConfiguration,omitempty"`
	DefaultRewriteRuleSet        *SubResource                            `json:"defaultRewriteRuleSet,omitempty"`
	ProvisioningState            *ProvisioningState                      `json:"provisioningState,omitempty"`
}
