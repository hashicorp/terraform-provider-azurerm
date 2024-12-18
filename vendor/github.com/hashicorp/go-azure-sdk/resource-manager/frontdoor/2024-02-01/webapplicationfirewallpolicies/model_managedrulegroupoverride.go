package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedRuleGroupOverride struct {
	Exclusions    *[]ManagedRuleExclusion `json:"exclusions,omitempty"`
	RuleGroupName string                  `json:"ruleGroupName"`
	Rules         *[]ManagedRuleOverride  `json:"rules,omitempty"`
}
