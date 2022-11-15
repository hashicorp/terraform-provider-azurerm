package webapplicationfirewallpolicies

type ManagedRuleExclusion struct {
	MatchVariable         ManagedRuleExclusionMatchVariable         `json:"matchVariable"`
	Selector              string                                    `json:"selector"`
	SelectorMatchOperator ManagedRuleExclusionSelectorMatchOperator `json:"selectorMatchOperator"`
}
