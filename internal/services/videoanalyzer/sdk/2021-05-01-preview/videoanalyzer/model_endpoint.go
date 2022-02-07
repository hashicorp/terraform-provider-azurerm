package videoanalyzer

type Endpoint struct {
	EndpointUrl *string                   `json:"endpointUrl,omitempty"`
	Type        VideoAnalyzerEndpointType `json:"type"`
}
