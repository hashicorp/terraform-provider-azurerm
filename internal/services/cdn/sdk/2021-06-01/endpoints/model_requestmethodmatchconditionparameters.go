package endpoints

type RequestMethodMatchConditionParameters struct {
	MatchValues     *[]MatchValues        `json:"matchValues,omitempty"`
	NegateCondition *bool                 `json:"negateCondition,omitempty"`
	Operator        RequestMethodOperator `json:"operator"`
	Transforms      *[]Transform          `json:"transforms,omitempty"`
	TypeName        TypeName              `json:"typeName"`
}
