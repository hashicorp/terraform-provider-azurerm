package frontdoors

type RulesEngineRule struct {
	Action                  RulesEngineAction            `json:"action"`
	MatchConditions         *[]RulesEngineMatchCondition `json:"matchConditions,omitempty"`
	MatchProcessingBehavior *MatchProcessingBehavior     `json:"matchProcessingBehavior,omitempty"`
	Name                    string                       `json:"name"`
	Priority                int64                        `json:"priority"`
}
