package endpoints

type SslProtocolMatchConditionParameters struct {
	MatchValues     *[]SslProtocol      `json:"matchValues,omitempty"`
	NegateCondition *bool               `json:"negateCondition,omitempty"`
	Operator        SslProtocolOperator `json:"operator"`
	Transforms      *[]Transform        `json:"transforms,omitempty"`
	TypeName        TypeName            `json:"typeName"`
}
