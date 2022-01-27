package webapplicationfirewallpolicies

type ManagedRuleSet struct {
	AnomalyScore       *int64                      `json:"anomalyScore,omitempty"`
	RuleGroupOverrides *[]ManagedRuleGroupOverride `json:"ruleGroupOverrides,omitempty"`
	RuleSetType        string                      `json:"ruleSetType"`
	RuleSetVersion     string                      `json:"ruleSetVersion"`
}
