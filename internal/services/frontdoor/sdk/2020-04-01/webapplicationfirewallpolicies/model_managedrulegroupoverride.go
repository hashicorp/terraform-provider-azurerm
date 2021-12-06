package webapplicationfirewallpolicies

type ManagedRuleGroupOverride struct {
	Exclusions    *[]ManagedRuleExclusion `json:"exclusions,omitempty"`
	RuleGroupName string                  `json:"ruleGroupName"`
	Rules         *[]ManagedRuleOverride  `json:"rules,omitempty"`
}
