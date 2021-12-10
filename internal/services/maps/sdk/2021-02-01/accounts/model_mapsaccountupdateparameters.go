package accounts

type MapsAccountUpdateParameters struct {
	Kind       *Kind                  `json:"kind,omitempty"`
	Properties *MapsAccountProperties `json:"properties,omitempty"`
	Sku        *Sku                   `json:"sku,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
