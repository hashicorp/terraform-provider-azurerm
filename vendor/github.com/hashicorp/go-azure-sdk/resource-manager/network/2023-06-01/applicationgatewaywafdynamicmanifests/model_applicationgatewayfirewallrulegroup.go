package applicationgatewaywafdynamicmanifests

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayFirewallRuleGroup struct {
	Description   *string                          `json:"description,omitempty"`
	RuleGroupName string                           `json:"ruleGroupName"`
	Rules         []ApplicationGatewayFirewallRule `json:"rules"`
}
