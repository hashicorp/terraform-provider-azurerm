package endpoints

type RequestSchemeMatchConditionParameters struct {
	MatchValues     *[]MatchValues `json:"matchValues,omitempty"`
	NegateCondition *bool          `json:"negateCondition,omitempty"`
	Operator        Operator       `json:"operator"`
	Transforms      *[]Transform   `json:"transforms,omitempty"`
	TypeName        TypeName       `json:"typeName"`
}
