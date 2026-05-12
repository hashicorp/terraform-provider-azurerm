package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExclusionManagedRuleGroup struct {
	RuleGroupName string                  `json:"ruleGroupName"`
	Rules         *[]ExclusionManagedRule `json:"rules,omitempty"`
}
