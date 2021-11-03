package cognitiveservicesaccounts

type Sku struct {
	Capacity *int64   `json:"capacity,omitempty"`
	Family   *string  `json:"family,omitempty"`
	Name     string   `json:"name"`
	Size     *string  `json:"size,omitempty"`
	Tier     *SkuTier `json:"tier,omitempty"`
}
