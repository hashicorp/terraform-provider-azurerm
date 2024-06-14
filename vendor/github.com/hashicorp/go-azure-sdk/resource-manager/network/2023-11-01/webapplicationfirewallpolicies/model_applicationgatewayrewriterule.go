package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRewriteRule struct {
	ActionSet    *ApplicationGatewayRewriteRuleActionSet   `json:"actionSet,omitempty"`
	Conditions   *[]ApplicationGatewayRewriteRuleCondition `json:"conditions,omitempty"`
	Name         *string                                   `json:"name,omitempty"`
	RuleSequence *int64                                    `json:"ruleSequence,omitempty"`
}
