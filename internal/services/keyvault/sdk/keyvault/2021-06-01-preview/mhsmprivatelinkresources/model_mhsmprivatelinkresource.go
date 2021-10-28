package mhsmprivatelinkresources

type MHSMPrivateLinkResource struct {
	Id         *string                            `json:"id,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Properties *MHSMPrivateLinkResourceProperties `json:"properties,omitempty"`
	Sku        *ManagedHsmSku                     `json:"sku,omitempty"`
	SystemData *SystemData                        `json:"systemData,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
