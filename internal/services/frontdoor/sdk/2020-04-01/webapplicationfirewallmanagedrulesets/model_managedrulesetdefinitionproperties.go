package webapplicationfirewallmanagedrulesets

type ManagedRuleSetDefinitionProperties struct {
	ProvisioningState *string                       `json:"provisioningState,omitempty"`
	RuleGroups        *[]ManagedRuleGroupDefinition `json:"ruleGroups,omitempty"`
	RuleSetId         *string                       `json:"ruleSetId,omitempty"`
	RuleSetType       *string                       `json:"ruleSetType,omitempty"`
	RuleSetVersion    *string                       `json:"ruleSetVersion,omitempty"`
}
