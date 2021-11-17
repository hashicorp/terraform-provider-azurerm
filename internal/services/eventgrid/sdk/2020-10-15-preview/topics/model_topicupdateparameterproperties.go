package topics

type TopicUpdateParameterProperties struct {
	InboundIpRules      *[]InboundIpRule     `json:"inboundIpRules,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
