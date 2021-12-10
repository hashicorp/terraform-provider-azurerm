package namespaces

type Sku struct {
	Name SkuName  `json:"name"`
	Tier *SkuTier `json:"tier,omitempty"`
}
