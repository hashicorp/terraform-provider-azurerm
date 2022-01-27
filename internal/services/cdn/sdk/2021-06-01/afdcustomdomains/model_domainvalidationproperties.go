package afdcustomdomains

type DomainValidationProperties struct {
	ExpirationDate  *string `json:"expirationDate,omitempty"`
	ValidationToken *string `json:"validationToken,omitempty"`
}
