package frontdoors

type FrontendEndpoint struct {
	Id         *string                     `json:"id,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties *FrontendEndpointProperties `json:"properties,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}
