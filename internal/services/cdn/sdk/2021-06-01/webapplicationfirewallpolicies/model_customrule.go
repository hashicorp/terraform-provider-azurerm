package webapplicationfirewallpolicies

type CustomRule struct {
	Action          ActionType              `json:"action"`
	EnabledState    *CustomRuleEnabledState `json:"enabledState,omitempty"`
	MatchConditions []MatchCondition        `json:"matchConditions"`
	Name            string                  `json:"name"`
	Priority        int64                   `json:"priority"`
}
