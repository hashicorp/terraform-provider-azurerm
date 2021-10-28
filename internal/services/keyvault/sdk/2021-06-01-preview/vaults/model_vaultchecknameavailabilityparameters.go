package vaults

type VaultCheckNameAvailabilityParameters struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}
