package managedapis

type ApiResourceGeneralInformation struct {
	Description   *string `json:"description,omitempty"`
	DisplayName   *string `json:"displayName,omitempty"`
	IconUrl       *string `json:"iconUrl,omitempty"`
	ReleaseTag    *string `json:"releaseTag,omitempty"`
	TermsOfUseUrl *string `json:"termsOfUseUrl,omitempty"`
}
