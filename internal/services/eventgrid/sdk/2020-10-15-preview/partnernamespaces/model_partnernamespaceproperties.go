package partnernamespaces

type PartnerNamespaceProperties struct {
	Endpoint                            *string                            `json:"endpoint,omitempty"`
	PartnerRegistrationFullyQualifiedId *string                            `json:"partnerRegistrationFullyQualifiedId,omitempty"`
	ProvisioningState                   *PartnerNamespaceProvisioningState `json:"provisioningState,omitempty"`
}
