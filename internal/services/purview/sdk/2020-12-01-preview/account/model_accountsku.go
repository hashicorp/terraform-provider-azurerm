package account

type AccountSku struct {
	Capacity *int64 `json:"capacity,omitempty"`
	Name     *Name  `json:"name,omitempty"`
}
