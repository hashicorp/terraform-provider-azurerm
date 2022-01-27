package endpoints

type UrlPathMatchConditionParameters struct {
	MatchValues     *[]string       `json:"matchValues,omitempty"`
	NegateCondition *bool           `json:"negateCondition,omitempty"`
	Operator        UrlPathOperator `json:"operator"`
	Transforms      *[]Transform    `json:"transforms,omitempty"`
	TypeName        TypeName        `json:"typeName"`
}
