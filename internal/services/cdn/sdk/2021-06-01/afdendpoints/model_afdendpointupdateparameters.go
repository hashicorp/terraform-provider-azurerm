package afdendpoints

type AFDEndpointUpdateParameters struct {
	Properties *AFDEndpointPropertiesUpdateParameters `json:"properties,omitempty"`
	Tags       *map[string]string                     `json:"tags,omitempty"`
}
