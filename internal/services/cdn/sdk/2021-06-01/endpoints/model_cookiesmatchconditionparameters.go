package endpoints

type CookiesMatchConditionParameters struct {
	MatchValues     *[]string       `json:"matchValues,omitempty"`
	NegateCondition *bool           `json:"negateCondition,omitempty"`
	Operator        CookiesOperator `json:"operator"`
	Selector        *string         `json:"selector,omitempty"`
	Transforms      *[]Transform    `json:"transforms,omitempty"`
	TypeName        TypeName        `json:"typeName"`
}
