package kubernetes

type CredentialResult struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}
