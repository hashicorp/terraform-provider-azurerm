package accounts

type MapsAccount struct {
	Id         *string                `json:"id,omitempty"`
	Kind       *Kind                  `json:"kind,omitempty"`
	Location   string                 `json:"location"`
	Name       *string                `json:"name,omitempty"`
	Properties *MapsAccountProperties `json:"properties,omitempty"`
	Sku        Sku                    `json:"sku"`
	SystemData *SystemData            `json:"systemData,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
