package servers

type ResourceSku struct {
	Capacity *int64   `json:"capacity,omitempty"`
	Name     string   `json:"name"`
	Tier     *SkuTier `json:"tier,omitempty"`
}
