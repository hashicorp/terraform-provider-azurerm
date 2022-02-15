package videoanalyzer

type KeyVaultProperties struct {
	CurrentKeyIdentifier *string `json:"currentKeyIdentifier,omitempty"`
	KeyIdentifier        string  `json:"keyIdentifier"`
}
