package secrets

type SecretProperties struct {
	Attributes           *Attributes `json:"attributes,omitempty"`
	ContentType          *string     `json:"contentType,omitempty"`
	SecretUri            *string     `json:"secretUri,omitempty"`
	SecretUriWithVersion *string     `json:"secretUriWithVersion,omitempty"`
	Value                *string     `json:"value,omitempty"`
}
