package cognitiveservicesaccounts

type KeyVaultProperties struct {
	IdentityClientId *string `json:"identityClientId,omitempty"`
	KeyName          *string `json:"keyName,omitempty"`
	KeyVaultUri      *string `json:"keyVaultUri,omitempty"`
	KeyVersion       *string `json:"keyVersion,omitempty"`
}
