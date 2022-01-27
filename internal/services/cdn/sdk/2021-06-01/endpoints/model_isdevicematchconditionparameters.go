package endpoints

type IsDeviceMatchConditionParameters struct {
	MatchValues     *[]MatchValues   `json:"matchValues,omitempty"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        IsDeviceOperator `json:"operator"`
	Transforms      *[]Transform     `json:"transforms,omitempty"`
	TypeName        TypeName         `json:"typeName"`
}
