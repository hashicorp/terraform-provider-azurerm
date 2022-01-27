package endpoints

type SocketAddrMatchConditionParameters struct {
	MatchValues     *[]string          `json:"matchValues,omitempty"`
	NegateCondition *bool              `json:"negateCondition,omitempty"`
	Operator        SocketAddrOperator `json:"operator"`
	Transforms      *[]Transform       `json:"transforms,omitempty"`
	TypeName        TypeName           `json:"typeName"`
}
