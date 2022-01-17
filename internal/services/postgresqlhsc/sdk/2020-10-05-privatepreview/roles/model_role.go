package roles

type Role struct {
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties *RoleProperties `json:"properties,omitempty"`
	SystemData *SystemData     `json:"systemData,omitempty"`
	Type       *string         `json:"type,omitempty"`
}
