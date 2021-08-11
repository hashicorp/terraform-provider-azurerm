package namespaces

type PrivateEndpointConnection struct {
	Id         *string                              `json:"id,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
	SystemData *SystemData                          `json:"systemData,omitempty"`
	Type       *string                              `json:"type,omitempty"`
}
