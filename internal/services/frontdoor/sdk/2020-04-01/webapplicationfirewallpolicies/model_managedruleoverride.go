package webapplicationfirewallpolicies

type ManagedRuleOverride struct {
	Action       *ActionType              `json:"action,omitempty"`
	EnabledState *ManagedRuleEnabledState `json:"enabledState,omitempty"`
	Exclusions   *[]ManagedRuleExclusion  `json:"exclusions,omitempty"`
	RuleId       string                   `json:"ruleId"`
}
