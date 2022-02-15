package accounts

type CreateTrustedIdProviderWithAccountParameters struct {
	Name       string                                    `json:"name"`
	Properties CreateOrUpdateTrustedIdProviderProperties `json:"properties"`
}
