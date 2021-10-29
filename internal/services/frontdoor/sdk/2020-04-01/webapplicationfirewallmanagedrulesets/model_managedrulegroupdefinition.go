package webapplicationfirewallmanagedrulesets

type ManagedRuleGroupDefinition struct {
	Description   *string                  `json:"description,omitempty"`
	RuleGroupName *string                  `json:"ruleGroupName,omitempty"`
	Rules         *[]ManagedRuleDefinition `json:"rules,omitempty"`
}
