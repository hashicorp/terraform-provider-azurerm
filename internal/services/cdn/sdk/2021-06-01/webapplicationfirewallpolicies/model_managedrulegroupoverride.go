package webapplicationfirewallpolicies

type ManagedRuleGroupOverride struct {
	RuleGroupName string                 `json:"ruleGroupName"`
	Rules         *[]ManagedRuleOverride `json:"rules,omitempty"`
}
