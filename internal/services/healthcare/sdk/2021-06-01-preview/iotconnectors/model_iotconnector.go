package iotconnectors

type IotConnector struct {
	Etag       *string                         `json:"etag,omitempty"`
	Id         *string                         `json:"id,omitempty"`
	Identity   *ServiceManagedIdentityIdentity `json:"identity,omitempty"`
	Location   *string                         `json:"location,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Properties *IotConnectorProperties         `json:"properties,omitempty"`
	SystemData *SystemData                     `json:"systemData,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}
