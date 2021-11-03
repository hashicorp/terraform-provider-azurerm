package accounts

type EncryptionConfig struct {
	KeyVaultMetaInfo *KeyVaultMetaInfo    `json:"keyVaultMetaInfo,omitempty"`
	Type             EncryptionConfigType `json:"type"`
}
