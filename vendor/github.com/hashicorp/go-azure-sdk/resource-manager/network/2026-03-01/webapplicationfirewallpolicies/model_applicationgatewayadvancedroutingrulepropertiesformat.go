package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAdvancedRoutingRulePropertiesFormat struct {
	AdvancedRoutingConditionSet *SubResource                    `json:"advancedRoutingConditionSet,omitempty"`
	AuthConfigs                 *[]ApplicationGatewayAuthConfig `json:"authConfigs,omitempty"`
	BackendAddressPool          *SubResource                    `json:"backendAddressPool,omitempty"`
	BackendHTTPSettings         *SubResource                    `json:"backendHttpSettings,omitempty"`
	Priority                    int64                           `json:"priority"`
	ProvisioningState           *ProvisioningState              `json:"provisioningState,omitempty"`
	RedirectConfiguration       *SubResource                    `json:"redirectConfiguration,omitempty"`
	RewriteRuleSet              *SubResource                    `json:"rewriteRuleSet,omitempty"`
}
