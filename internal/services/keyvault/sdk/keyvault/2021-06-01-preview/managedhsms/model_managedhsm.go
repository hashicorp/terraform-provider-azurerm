package managedhsms

type ManagedHsm struct {
	Id         *string               `json:"id,omitempty"`
	Location   *string               `json:"location,omitempty"`
	Name       *string               `json:"name,omitempty"`
	Properties *ManagedHsmProperties `json:"properties,omitempty"`
	Sku        *ManagedHsmSku        `json:"sku,omitempty"`
	SystemData *SystemData           `json:"systemData,omitempty"`
	Tags       *map[string]string    `json:"tags,omitempty"`
	Type       *string               `json:"type,omitempty"`
}
