package resource

type ServicesPatchDescription struct {
	Properties *ServicesPropertiesUpdateParameters `json:"properties,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
}
