package endpoints

type Endpoint struct {
	Id         *string             `json:"id,omitempty"`
	Location   string              `json:"location"`
	Name       *string             `json:"name,omitempty"`
	Properties *EndpointProperties `json:"properties,omitempty"`
	SystemData *SystemData         `json:"systemData,omitempty"`
	Tags       *map[string]string  `json:"tags,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
