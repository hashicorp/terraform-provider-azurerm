package signalr

type UpstreamTemplate struct {
	Auth            *UpstreamAuthSettings `json:"auth,omitempty"`
	CategoryPattern *string               `json:"categoryPattern,omitempty"`
	EventPattern    *string               `json:"eventPattern,omitempty"`
	HubPattern      *string               `json:"hubPattern,omitempty"`
	UrlTemplate     string                `json:"urlTemplate"`
}
