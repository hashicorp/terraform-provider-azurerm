package diskpoolzones

type Sku struct {
	Name string  `json:"name"`
	Tier *string `json:"tier,omitempty"`
}
