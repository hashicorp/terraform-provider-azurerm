package configurationstores

type KeyVaultProperties struct {
	IdentityClientId *string `json:"identityClientId,omitempty"`
	KeyIdentifier    *string `json:"keyIdentifier,omitempty"`
}
