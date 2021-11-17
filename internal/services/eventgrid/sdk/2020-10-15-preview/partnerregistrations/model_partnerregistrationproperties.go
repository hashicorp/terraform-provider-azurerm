package partnerregistrations

type PartnerRegistrationProperties struct {
	AuthorizedAzureSubscriptionIds  *[]string                             `json:"authorizedAzureSubscriptionIds,omitempty"`
	CustomerServiceUri              *string                               `json:"customerServiceUri,omitempty"`
	LogoUri                         *string                               `json:"logoUri,omitempty"`
	LongDescription                 *string                               `json:"longDescription,omitempty"`
	PartnerCustomerServiceExtension *string                               `json:"partnerCustomerServiceExtension,omitempty"`
	PartnerCustomerServiceNumber    *string                               `json:"partnerCustomerServiceNumber,omitempty"`
	PartnerName                     *string                               `json:"partnerName,omitempty"`
	PartnerResourceTypeDescription  *string                               `json:"partnerResourceTypeDescription,omitempty"`
	PartnerResourceTypeDisplayName  *string                               `json:"partnerResourceTypeDisplayName,omitempty"`
	PartnerResourceTypeName         *string                               `json:"partnerResourceTypeName,omitempty"`
	ProvisioningState               *PartnerRegistrationProvisioningState `json:"provisioningState,omitempty"`
	SetupUri                        *string                               `json:"setupUri,omitempty"`
	VisibilityState                 *PartnerRegistrationVisibilityState   `json:"visibilityState,omitempty"`
}
