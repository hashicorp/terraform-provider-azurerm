package secrets

type SecretPatchParameters struct {
	Properties *SecretPatchProperties `json:"properties,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
