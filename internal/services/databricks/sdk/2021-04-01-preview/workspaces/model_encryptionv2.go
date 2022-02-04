package workspaces

type EncryptionV2 struct {
	KeySource          EncryptionKeySource             `json:"keySource"`
	KeyVaultProperties *EncryptionV2KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
