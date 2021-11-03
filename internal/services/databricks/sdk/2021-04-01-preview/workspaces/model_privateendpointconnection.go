package workspaces

type PrivateEndpointConnection struct {
	Id         *string                             `json:"id,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties PrivateEndpointConnectionProperties `json:"properties"`
	Type       *string                             `json:"type,omitempty"`
}
