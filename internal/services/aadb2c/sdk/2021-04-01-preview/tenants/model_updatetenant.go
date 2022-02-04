package tenants

type UpdateTenant struct {
	Properties UpdateTenantProperties `json:"properties"`
	Sku        Sku                    `json:"sku"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
