package account

type AccountEndpoints struct {
	Catalog  *string `json:"catalog,omitempty"`
	Guardian *string `json:"guardian,omitempty"`
	Scan     *string `json:"scan,omitempty"`
}
