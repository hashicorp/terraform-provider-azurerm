package accounts

type UpdateTrustedIdProviderWithAccountParameters struct {
	Name       string                             `json:"name"`
	Properties *UpdateTrustedIdProviderProperties `json:"properties,omitempty"`
}
