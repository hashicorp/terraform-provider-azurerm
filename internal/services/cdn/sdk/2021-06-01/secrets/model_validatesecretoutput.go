package secrets

type ValidateSecretOutput struct {
	Message *string `json:"message,omitempty"`
	Status  *Status `json:"status,omitempty"`
}
