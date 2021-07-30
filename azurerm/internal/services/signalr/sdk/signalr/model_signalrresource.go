package signalr

type SignalRResource struct {
	Id         *string            `json:"id,omitempty"`
	Kind       *ServiceKind       `json:"kind,omitempty"`
	Location   *string            `json:"location,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *SignalRProperties `json:"properties,omitempty"`
	Sku        *ResourceSku       `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
