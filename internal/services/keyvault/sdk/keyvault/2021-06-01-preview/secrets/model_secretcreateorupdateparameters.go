package secrets

type SecretCreateOrUpdateParameters struct {
	Properties SecretProperties   `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
