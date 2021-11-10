package ipfilterrules

type IpFilterRuleProperties struct {
	Action     *IPAction `json:"action,omitempty"`
	FilterName *string   `json:"filterName,omitempty"`
	IpMask     *string   `json:"ipMask,omitempty"`
}
