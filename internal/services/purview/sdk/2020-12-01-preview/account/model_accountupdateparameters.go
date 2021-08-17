package account

type AccountUpdateParameters struct {
	Properties *AccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
