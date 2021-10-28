package vaults

type VaultPatchParameters struct {
	Properties *VaultPatchProperties `json:"properties,omitempty"`
	Tags       *map[string]string    `json:"tags,omitempty"`
}
