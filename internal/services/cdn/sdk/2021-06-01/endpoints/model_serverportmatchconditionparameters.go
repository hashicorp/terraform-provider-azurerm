package endpoints

type ServerPortMatchConditionParameters struct {
	MatchValues     *[]string          `json:"matchValues,omitempty"`
	NegateCondition *bool              `json:"negateCondition,omitempty"`
	Operator        ServerPortOperator `json:"operator"`
	Transforms      *[]Transform       `json:"transforms,omitempty"`
	TypeName        TypeName           `json:"typeName"`
}
