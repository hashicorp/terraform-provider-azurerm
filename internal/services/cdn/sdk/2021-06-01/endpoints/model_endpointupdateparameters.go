package endpoints

type EndpointUpdateParameters struct {
	Properties *EndpointPropertiesUpdateParameters `json:"properties,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
}
