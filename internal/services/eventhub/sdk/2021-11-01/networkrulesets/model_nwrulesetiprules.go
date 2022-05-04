package networkrulesets

type NWRuleSetIpRules struct {
	Action *NetworkRuleIPAction `json:"action,omitempty"`
	IpMask *string              `json:"ipMask,omitempty"`
}
