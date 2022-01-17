package servers

type ServerGroupServer struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *ServerGroupServerProperties `json:"properties,omitempty"`
	SystemData *SystemData                  `json:"systemData,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
