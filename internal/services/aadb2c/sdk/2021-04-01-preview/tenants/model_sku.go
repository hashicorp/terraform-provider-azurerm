package tenants

type Sku struct {
	Name SkuName `json:"name"`
	Tier SkuTier `json:"tier"`
}
