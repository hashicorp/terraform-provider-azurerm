package account

type AccountPropertiesEndpoints struct {
	Catalog  *string `json:"catalog,omitempty"`
	Guardian *string `json:"guardian,omitempty"`
	Scan     *string `json:"scan,omitempty"`
}
