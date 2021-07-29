package servers

type SkuDetailsForExistingResource struct {
	ResourceType *string      `json:"resourceType,omitempty"`
	Sku          *ResourceSku `json:"sku,omitempty"`
}
