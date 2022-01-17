package configurations

type ServerGroupConfiguration struct {
	Id         *string                             `json:"id,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *ServerGroupConfigurationProperties `json:"properties,omitempty"`
	SystemData *SystemData                         `json:"systemData,omitempty"`
	Type       *string                             `json:"type,omitempty"`
}
