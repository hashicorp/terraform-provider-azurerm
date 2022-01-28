package fhirservices

type FhirServiceAuthenticationConfiguration struct {
	Audience          *string `json:"audience,omitempty"`
	Authority         *string `json:"authority,omitempty"`
	SmartProxyEnabled *bool   `json:"smartProxyEnabled,omitempty"`
}
