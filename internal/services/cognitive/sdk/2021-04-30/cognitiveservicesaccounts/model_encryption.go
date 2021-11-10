package cognitiveservicesaccounts

type Encryption struct {
	KeySource          *KeySource          `json:"keySource,omitempty"`
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
}
