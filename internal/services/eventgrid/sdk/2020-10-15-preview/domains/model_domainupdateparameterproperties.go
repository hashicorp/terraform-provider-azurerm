package domains

type DomainUpdateParameterProperties struct {
	InboundIpRules      *[]InboundIpRule     `json:"inboundIpRules,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
