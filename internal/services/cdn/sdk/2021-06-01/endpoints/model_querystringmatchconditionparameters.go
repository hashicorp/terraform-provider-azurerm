package endpoints

type QueryStringMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        QueryStringOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`
	TypeName        TypeName            `json:"typeName"`
}
