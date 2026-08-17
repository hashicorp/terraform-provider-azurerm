package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRoutingRulePropertiesFormat struct {
	BackendAddressPool *SubResource                              `json:"backendAddressPool,omitempty"`
	BackendSettings    *SubResource                              `json:"backendSettings,omitempty"`
	Listener           *SubResource                              `json:"listener,omitempty"`
	Priority           int64                                     `json:"priority"`
	ProvisioningState  *ProvisioningState                        `json:"provisioningState,omitempty"`
	RuleType           *ApplicationGatewayRequestRoutingRuleType `json:"ruleType,omitempty"`
}
