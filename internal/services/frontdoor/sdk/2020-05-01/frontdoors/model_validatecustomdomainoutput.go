package frontdoors

type ValidateCustomDomainOutput struct {
	CustomDomainValidated *bool   `json:"customDomainValidated,omitempty"`
	Message               *string `json:"message,omitempty"`
	Reason                *string `json:"reason,omitempty"`
}
