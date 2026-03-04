package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayFirewallRuleSetPropertiesFormat struct {
	ProvisioningState *ProvisioningState                    `json:"provisioningState,omitempty"`
	RuleGroups        []ApplicationGatewayFirewallRuleGroup `json:"ruleGroups"`
	RuleSetType       string                                `json:"ruleSetType"`
	RuleSetVersion    string                                `json:"ruleSetVersion"`
	Tiers             *[]ApplicationGatewayTierTypes        `json:"tiers,omitempty"`
}
