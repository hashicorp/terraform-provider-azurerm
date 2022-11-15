package webapplicationfirewallpolicies

type ManagedRuleSet struct {
	Exclusions         *[]ManagedRuleExclusion     `json:"exclusions,omitempty"`
	RuleGroupOverrides *[]ManagedRuleGroupOverride `json:"ruleGroupOverrides,omitempty"`
	RuleSetType        string                      `json:"ruleSetType"`
	RuleSetVersion     string                      `json:"ruleSetVersion"`
}
