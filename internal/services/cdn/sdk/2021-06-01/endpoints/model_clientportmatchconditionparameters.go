package endpoints

type ClientPortMatchConditionParameters struct {
	MatchValues     *[]string          `json:"matchValues,omitempty"`
	NegateCondition *bool              `json:"negateCondition,omitempty"`
	Operator        ClientPortOperator `json:"operator"`
	Transforms      *[]Transform       `json:"transforms,omitempty"`
	TypeName        TypeName           `json:"typeName"`
}
