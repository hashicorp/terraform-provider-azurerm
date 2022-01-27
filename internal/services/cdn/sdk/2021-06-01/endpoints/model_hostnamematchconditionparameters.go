package endpoints

type HostNameMatchConditionParameters struct {
	MatchValues     *[]string        `json:"matchValues,omitempty"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        HostNameOperator `json:"operator"`
	Transforms      *[]Transform     `json:"transforms,omitempty"`
	TypeName        TypeName         `json:"typeName"`
}
