package videoanalyzer

type VideoEntity struct {
	Id         *string          `json:"id,omitempty"`
	Name       *string          `json:"name,omitempty"`
	Properties *VideoProperties `json:"properties,omitempty"`
	SystemData *SystemData      `json:"systemData,omitempty"`
	Type       *string          `json:"type,omitempty"`
}
