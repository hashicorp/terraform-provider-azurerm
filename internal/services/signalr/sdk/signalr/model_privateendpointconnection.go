package signalr

type PrivateEndpointConnection struct {
	Id         *string                              `json:"id,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
	Type       *string                              `json:"type,omitempty"`
}
