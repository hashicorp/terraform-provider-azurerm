package endpoints

type RequestUriMatchConditionParameters struct {
	MatchValues     *[]string          `json:"matchValues,omitempty"`
	NegateCondition *bool              `json:"negateCondition,omitempty"`
	Operator        RequestUriOperator `json:"operator"`
	Transforms      *[]Transform       `json:"transforms,omitempty"`
	TypeName        TypeName           `json:"typeName"`
}
