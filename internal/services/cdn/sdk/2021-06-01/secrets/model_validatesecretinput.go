package secrets

type ValidateSecretInput struct {
	SecretSource  ResourceReference `json:"secretSource"`
	SecretType    SecretType        `json:"secretType"`
	SecretVersion *string           `json:"secretVersion,omitempty"`
}
