package secrets

type SecretPatchProperties struct {
	Attributes  *Attributes `json:"attributes,omitempty"`
	ContentType *string     `json:"contentType,omitempty"`
	Value       *string     `json:"value,omitempty"`
}
