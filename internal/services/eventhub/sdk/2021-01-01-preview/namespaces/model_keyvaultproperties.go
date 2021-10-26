package namespaces

type KeyVaultProperties struct {
	Identity    *UserAssignedIdentityProperties `json:"identity,omitempty"`
	KeyName     *string                         `json:"keyName,omitempty"`
	KeyVaultUri *string                         `json:"keyVaultUri,omitempty"`
	KeyVersion  *string                         `json:"keyVersion,omitempty"`
}
