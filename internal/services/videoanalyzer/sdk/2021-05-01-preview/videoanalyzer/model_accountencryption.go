package videoanalyzer

type AccountEncryption struct {
	Identity           *ResourceIdentity        `json:"identity,omitempty"`
	KeyVaultProperties *KeyVaultProperties      `json:"keyVaultProperties,omitempty"`
	Status             *string                  `json:"status,omitempty"`
	Type               AccountEncryptionKeyType `json:"type"`
}
