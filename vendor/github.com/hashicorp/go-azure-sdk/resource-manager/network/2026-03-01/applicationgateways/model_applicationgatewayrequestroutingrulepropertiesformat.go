package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRequestRoutingRulePropertiesFormat struct {
	AdvancedRoutingMap       *SubResource                              `json:"advancedRoutingMap,omitempty"`
	AuthConfigs              *[]ApplicationGatewayAuthConfig           `json:"authConfigs,omitempty"`
	BackendAddressPool       *SubResource                              `json:"backendAddressPool,omitempty"`
	BackendHTTPSettings      *SubResource                              `json:"backendHttpSettings,omitempty"`
	EntraJWTValidationConfig *SubResource                              `json:"entraJWTValidationConfig,omitempty"`
	HTTPListener             *SubResource                              `json:"httpListener,omitempty"`
	LoadDistributionPolicy   *SubResource                              `json:"loadDistributionPolicy,omitempty"`
	Priority                 *int64                                    `json:"priority,omitempty"`
	ProvisioningState        *ProvisioningState                        `json:"provisioningState,omitempty"`
	RedirectConfiguration    *SubResource                              `json:"redirectConfiguration,omitempty"`
	RewriteRuleSet           *SubResource                              `json:"rewriteRuleSet,omitempty"`
	RuleType                 *ApplicationGatewayRequestRoutingRuleType `json:"ruleType,omitempty"`
	UrlPathMap               *SubResource                              `json:"urlPathMap,omitempty"`
}
