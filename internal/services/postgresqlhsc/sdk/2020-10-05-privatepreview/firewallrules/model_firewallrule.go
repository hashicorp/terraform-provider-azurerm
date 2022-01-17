package firewallrules

type FirewallRule struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties FirewallRuleProperties `json:"properties"`
	SystemData *SystemData            `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
