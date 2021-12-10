package namespaces

type Encryption struct {
	KeySource                       *KeySource            `json:"keySource,omitempty"`
	KeyVaultProperties              *[]KeyVaultProperties `json:"keyVaultProperties,omitempty"`
	RequireInfrastructureEncryption *bool                 `json:"requireInfrastructureEncryption,omitempty"`
}
