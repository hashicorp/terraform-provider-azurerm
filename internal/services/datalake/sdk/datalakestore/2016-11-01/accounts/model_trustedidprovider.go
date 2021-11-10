package accounts

type TrustedIdProvider struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *TrustedIdProviderProperties `json:"properties,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
