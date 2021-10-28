package frontdoors

type RulesEngineMatchCondition struct {
	NegateCondition          *bool                    `json:"negateCondition,omitempty"`
	RulesEngineMatchValue    []string                 `json:"rulesEngineMatchValue"`
	RulesEngineMatchVariable RulesEngineMatchVariable `json:"rulesEngineMatchVariable"`
	RulesEngineOperator      RulesEngineOperator      `json:"rulesEngineOperator"`
	Selector                 *string                  `json:"selector,omitempty"`
	Transforms               *[]Transform             `json:"transforms,omitempty"`
}
