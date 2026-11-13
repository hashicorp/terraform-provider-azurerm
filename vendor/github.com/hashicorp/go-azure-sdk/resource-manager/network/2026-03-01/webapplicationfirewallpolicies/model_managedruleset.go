package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedRuleSet struct {
	ComputedDisabledRules *[]ManagedRuleSetRuleGroup  `json:"computedDisabledRules,omitempty"`
	RuleGroupOverrides    *[]ManagedRuleGroupOverride `json:"ruleGroupOverrides,omitempty"`
	RuleSetType           string                      `json:"ruleSetType"`
	RuleSetVersion        string                      `json:"ruleSetVersion"`
}
