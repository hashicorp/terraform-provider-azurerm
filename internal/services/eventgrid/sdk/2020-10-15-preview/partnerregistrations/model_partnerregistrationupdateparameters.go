package partnerregistrations

type PartnerRegistrationUpdateParameters struct {
	AuthorizedAzureSubscriptionIds *[]string          `json:"authorizedAzureSubscriptionIds,omitempty"`
	LogoUri                        *string            `json:"logoUri,omitempty"`
	PartnerTopicTypeDescription    *string            `json:"partnerTopicTypeDescription,omitempty"`
	PartnerTopicTypeDisplayName    *string            `json:"partnerTopicTypeDisplayName,omitempty"`
	PartnerTopicTypeName           *string            `json:"partnerTopicTypeName,omitempty"`
	SetupUri                       *string            `json:"setupUri,omitempty"`
	Tags                           *map[string]string `json:"tags,omitempty"`
}
