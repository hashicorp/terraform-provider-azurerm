package rules

type RuleProperties struct {
	Actions                 []DeliveryRuleAction     `json:"actions"`
	Conditions              *[]DeliveryRuleCondition `json:"conditions,omitempty"`
	DeploymentStatus        *DeploymentStatus        `json:"deploymentStatus,omitempty"`
	MatchProcessingBehavior *MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
	Order                   int64                    `json:"order"`
	ProvisioningState       *AfdProvisioningState    `json:"provisioningState,omitempty"`
	RuleSetName             *string                  `json:"ruleSetName,omitempty"`
}
