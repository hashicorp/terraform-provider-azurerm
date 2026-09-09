package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayWebApplicationFirewallConfiguration struct {
	DisabledRuleGroups     *[]ApplicationGatewayFirewallDisabledRuleGroup `json:"disabledRuleGroups,omitempty"`
	Enabled                bool                                           `json:"enabled"`
	Exclusions             *[]ApplicationGatewayFirewallExclusion         `json:"exclusions,omitempty"`
	FileUploadLimitInMb    *int64                                         `json:"fileUploadLimitInMb,omitempty"`
	FirewallMode           ApplicationGatewayFirewallMode                 `json:"firewallMode"`
	MaxRequestBodySize     *int64                                         `json:"maxRequestBodySize,omitempty"`
	MaxRequestBodySizeInKb *int64                                         `json:"maxRequestBodySizeInKb,omitempty"`
	RequestBodyCheck       *bool                                          `json:"requestBodyCheck,omitempty"`
	RuleSetType            string                                         `json:"ruleSetType"`
	RuleSetVersion         string                                         `json:"ruleSetVersion"`
}
