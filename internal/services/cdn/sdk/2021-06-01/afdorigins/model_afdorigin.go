package afdorigins

type AFDOrigin struct {
	Id         *string              `json:"id,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Properties *AFDOriginProperties `json:"properties,omitempty"`
	SystemData *SystemData          `json:"systemData,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
