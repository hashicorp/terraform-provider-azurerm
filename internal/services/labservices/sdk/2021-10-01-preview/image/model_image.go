package image

type Image struct {
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties ImageProperties `json:"properties"`
	SystemData *SystemData     `json:"systemData,omitempty"`
	Type       *string         `json:"type,omitempty"`
}
