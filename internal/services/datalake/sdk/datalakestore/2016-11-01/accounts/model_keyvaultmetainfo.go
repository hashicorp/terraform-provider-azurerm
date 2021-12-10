package accounts

type KeyVaultMetaInfo struct {
	EncryptionKeyName    string `json:"encryptionKeyName"`
	EncryptionKeyVersion string `json:"encryptionKeyVersion"`
	KeyVaultResourceId   string `json:"keyVaultResourceId"`
}
