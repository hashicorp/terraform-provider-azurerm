package configurationstores

type EncryptionProperties struct {
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
