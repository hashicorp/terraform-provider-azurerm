package endpoints

type PostArgsMatchConditionParameters struct {
	MatchValues     *[]string        `json:"matchValues,omitempty"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        PostArgsOperator `json:"operator"`
	Selector        *string          `json:"selector,omitempty"`
	Transforms      *[]Transform     `json:"transforms,omitempty"`
	TypeName        TypeName         `json:"typeName"`
}
