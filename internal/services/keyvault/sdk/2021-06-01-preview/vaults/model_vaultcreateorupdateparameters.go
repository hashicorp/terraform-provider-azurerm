package vaults

type VaultCreateOrUpdateParameters struct {
	Location   string             `json:"location"`
	Properties VaultProperties    `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
