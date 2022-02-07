package workspaces

type EncryptionV2KeyVaultProperties struct {
	KeyName     string `json:"keyName"`
	KeyVaultUri string `json:"keyVaultUri"`
	KeyVersion  string `json:"keyVersion"`
}
