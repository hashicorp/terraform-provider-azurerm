package signalr

type UpstreamTemplate struct {
	CategoryPattern *string `json:"categoryPattern,omitempty"`
	EventPattern    *string `json:"eventPattern,omitempty"`
	HubPattern      *string `json:"hubPattern,omitempty"`
	UrlTemplate     string  `json:"urlTemplate"`
}
