package endpoints

type UrlFileNameMatchConditionParameters struct {
	MatchValues     *[]string           `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        UrlFileNameOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`
	TypeName        TypeName            `json:"typeName"`
}
