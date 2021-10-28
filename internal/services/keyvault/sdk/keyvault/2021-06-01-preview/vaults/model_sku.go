package vaults

type Sku struct {
	Family SkuFamily `json:"family"`
	Name   SkuName   `json:"name"`
}
