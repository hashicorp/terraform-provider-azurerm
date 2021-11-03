package videoanalyzer

type EdgeModuleEntity struct {
	Id         *string               `json:"id,omitempty"`
	Name       *string               `json:"name,omitempty"`
	Properties *EdgeModuleProperties `json:"properties,omitempty"`
	SystemData *SystemData           `json:"systemData,omitempty"`
	Type       *string               `json:"type,omitempty"`
}
