package namespaces

type KeyVaultProperties struct {
	KeyName     *string `json:"keyName,omitempty"`
	KeyVaultUri *string `json:"keyVaultUri,omitempty"`
	KeyVersion  *string `json:"keyVersion,omitempty"`
}
