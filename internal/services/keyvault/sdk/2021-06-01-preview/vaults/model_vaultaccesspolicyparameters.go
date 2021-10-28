package vaults

type VaultAccessPolicyParameters struct {
	Id         *string                     `json:"id,omitempty"`
	Location   *string                     `json:"location,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties VaultAccessPolicyProperties `json:"properties"`
	Type       *string                     `json:"type,omitempty"`
}
