package rules

type RuleUpdatePropertiesParameters struct {
	Actions                 *[]DeliveryRuleAction    `json:"actions,omitempty"`
	Conditions              *[]DeliveryRuleCondition `json:"conditions,omitempty"`
	MatchProcessingBehavior *MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
	Order                   *int64                   `json:"order,omitempty"`
	RuleSetName             *string                  `json:"ruleSetName,omitempty"`
}
