package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedRuleExclusion struct {
	MatchVariable         ManagedRuleExclusionMatchVariable         `json:"matchVariable"`
	Selector              string                                    `json:"selector"`
	SelectorMatchOperator ManagedRuleExclusionSelectorMatchOperator `json:"selectorMatchOperator"`
}
