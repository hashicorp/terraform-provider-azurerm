package endpoints

type Endpoint struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *EndpointProperties `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
