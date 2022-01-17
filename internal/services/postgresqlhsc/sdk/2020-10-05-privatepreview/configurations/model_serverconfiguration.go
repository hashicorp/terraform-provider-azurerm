package configurations

type ServerConfiguration struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *ServerConfigurationProperties `json:"properties,omitempty"`
	SystemData *SystemData                    `json:"systemData,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
