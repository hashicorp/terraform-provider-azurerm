package signalr

type Sku struct {
	Capacity     *SkuCapacity `json:"capacity,omitempty"`
	ResourceType *string      `json:"resourceType,omitempty"`
	Sku          *ResourceSku `json:"sku,omitempty"`
}
