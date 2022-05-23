package links

type LinkerResource struct {
	Id         *string          `json:"id,omitempty"`
	Name       *string          `json:"name,omitempty"`
	Properties LinkerProperties `json:"properties"`
	SystemData *SystemData      `json:"systemData,omitempty"`
	Type       *string          `json:"type,omitempty"`
}
