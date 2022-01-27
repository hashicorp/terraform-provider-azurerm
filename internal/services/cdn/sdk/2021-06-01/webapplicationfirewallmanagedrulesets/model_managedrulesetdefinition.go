package webapplicationfirewallmanagedrulesets

type ManagedRuleSetDefinition struct {
	Id         *string                             `json:"id,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *ManagedRuleSetDefinitionProperties `json:"properties,omitempty"`
	Sku        *Sku                                `json:"sku,omitempty"`
	SystemData *SystemData                         `json:"systemData,omitempty"`
	Type       *string                             `json:"type,omitempty"`
}
