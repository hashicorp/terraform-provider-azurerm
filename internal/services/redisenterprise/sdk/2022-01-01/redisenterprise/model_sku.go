package redisenterprise

type Sku struct {
	Capacity *int64  `json:"capacity,omitempty"`
	Name     SkuName `json:"name"`
}
