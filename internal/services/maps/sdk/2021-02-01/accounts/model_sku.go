package accounts

type Sku struct {
	Name Name    `json:"name"`
	Tier *string `json:"tier,omitempty"`
}
