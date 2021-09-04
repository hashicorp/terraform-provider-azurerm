package namespaces

type Sku struct {
	Capacity *int64   `json:"capacity,omitempty"`
	Name     SkuName  `json:"name"`
	Tier     *SkuTier `json:"tier,omitempty"`
}
