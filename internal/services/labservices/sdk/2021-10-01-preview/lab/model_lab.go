package lab

type Lab struct {
	Id         *string            `json:"id,omitempty"`
	Location   string             `json:"location"`
	Name       *string            `json:"name,omitempty"`
	Properties LabProperties      `json:"properties"`
	SystemData *SystemData        `json:"systemData,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
