package vaults

type DeletedVault struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *DeletedVaultProperties `json:"properties,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
