package endpoints

type UrlFileExtensionMatchConditionParameters struct {
	MatchValues     *[]string                `json:"matchValues,omitempty"`
	NegateCondition *bool                    `json:"negateCondition,omitempty"`
	Operator        UrlFileExtensionOperator `json:"operator"`
	Transforms      *[]Transform             `json:"transforms,omitempty"`
	TypeName        TypeName                 `json:"typeName"`
}
