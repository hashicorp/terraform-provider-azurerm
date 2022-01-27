package endpoints

type RequestBodyMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        RequestBodyOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`
	TypeName        TypeName            `json:"typeName"`
}
