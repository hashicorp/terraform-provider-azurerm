package servers

type IPv4FirewallRule struct {
	FirewallRuleName *string `json:"firewallRuleName,omitempty"`
	RangeEnd         *string `json:"rangeEnd,omitempty"`
	RangeStart       *string `json:"rangeStart,omitempty"`
}
