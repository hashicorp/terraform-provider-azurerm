package endpoints

type HttpVersionMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        HttpVersionOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`
	TypeName        TypeName            `json:"typeName"`
}
