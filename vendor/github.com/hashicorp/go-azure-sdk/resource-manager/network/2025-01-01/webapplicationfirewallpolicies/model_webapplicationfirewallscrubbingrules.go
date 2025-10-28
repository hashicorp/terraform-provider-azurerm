package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallScrubbingRules struct {
	MatchVariable         ScrubbingRuleEntryMatchVariable `json:"matchVariable"`
	Selector              *string                         `json:"selector,omitempty"`
	SelectorMatchOperator ScrubbingRuleEntryMatchOperator `json:"selectorMatchOperator"`
	State                 *ScrubbingRuleEntryState        `json:"state,omitempty"`
}
