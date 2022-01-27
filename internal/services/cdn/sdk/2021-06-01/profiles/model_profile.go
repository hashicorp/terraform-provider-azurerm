package profiles

type Profile struct {
	Id         *string            `json:"id,omitempty"`
	Kind       *string            `json:"kind,omitempty"`
	Location   string             `json:"location"`
	Name       *string            `json:"name,omitempty"`
	Properties *ProfileProperties `json:"properties,omitempty"`
	Sku        Sku                `json:"sku"`
	SystemData *SystemData        `json:"systemData,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
