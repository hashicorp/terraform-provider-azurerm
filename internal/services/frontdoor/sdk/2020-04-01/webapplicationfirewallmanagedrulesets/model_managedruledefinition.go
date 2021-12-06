package webapplicationfirewallmanagedrulesets

type ManagedRuleDefinition struct {
	DefaultAction *ActionType              `json:"defaultAction,omitempty"`
	DefaultState  *ManagedRuleEnabledState `json:"defaultState,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	RuleId        *string                  `json:"ruleId,omitempty"`
}
