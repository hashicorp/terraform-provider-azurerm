package endpoints

type RemoteAddressMatchConditionParameters struct {
	MatchValues     *[]string             `json:"matchValues,omitempty"`
	NegateCondition *bool                 `json:"negateCondition,omitempty"`
	Operator        RemoteAddressOperator `json:"operator"`
	Transforms      *[]Transform          `json:"transforms,omitempty"`
	TypeName        TypeName              `json:"typeName"`
}
