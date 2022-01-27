package origingroups

type OriginGroup struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *OriginGroupProperties `json:"properties,omitempty"`
	SystemData *SystemData            `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
