package edgenodes

type EdgeNode struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *EdgeNodeProperties `json:"properties,omitempty"`
	SystemData *SystemData         `json:"systemData,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
