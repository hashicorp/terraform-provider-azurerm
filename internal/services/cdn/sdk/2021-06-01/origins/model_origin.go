package origins

type Origin struct {
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties *OriginProperties `json:"properties,omitempty"`
	SystemData *SystemData       `json:"systemData,omitempty"`
	Type       *string           `json:"type,omitempty"`
}
