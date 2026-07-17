package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OwaspCrsExclusionEntry struct {
	ExclusionManagedRuleSets *[]ExclusionManagedRuleSet                  `json:"exclusionManagedRuleSets,omitempty"`
	MatchVariable            OwaspCrsExclusionEntryMatchVariable         `json:"matchVariable"`
	Selector                 string                                      `json:"selector"`
	SelectorMatchOperator    OwaspCrsExclusionEntrySelectorMatchOperator `json:"selectorMatchOperator"`
}
