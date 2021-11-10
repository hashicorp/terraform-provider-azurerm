package networkrulesets

type NWRuleSetVirtualNetworkRules struct {
	IgnoreMissingVnetServiceEndpoint *bool   `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
	Subnet                           *Subnet `json:"subnet,omitempty"`
}
