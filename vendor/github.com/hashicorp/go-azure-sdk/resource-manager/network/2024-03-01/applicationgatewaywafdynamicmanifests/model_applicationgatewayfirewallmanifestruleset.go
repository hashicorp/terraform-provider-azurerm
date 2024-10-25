package applicationgatewaywafdynamicmanifests

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayFirewallManifestRuleSet struct {
	RuleGroups     []ApplicationGatewayFirewallRuleGroup   `json:"ruleGroups"`
	RuleSetType    string                                  `json:"ruleSetType"`
	RuleSetVersion string                                  `json:"ruleSetVersion"`
	Status         *ApplicationGatewayRuleSetStatusOptions `json:"status,omitempty"`
	Tiers          *[]ApplicationGatewayTierTypes          `json:"tiers,omitempty"`
}
