package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExceptionEntry struct {
	ExceptionManagedRuleSets *[]ExclusionManagedRuleSet           `json:"exceptionManagedRuleSets,omitempty"`
	MatchVariable            ExceptionEntryMatchVariable          `json:"matchVariable"`
	Selector                 *string                              `json:"selector,omitempty"`
	SelectorMatchOperator    *ExceptionEntrySelectorMatchOperator `json:"selectorMatchOperator,omitempty"`
	ValueMatchOperator       ExceptionEntryValueMatchOperator     `json:"valueMatchOperator"`
	Values                   *[]string                            `json:"values,omitempty"`
}
