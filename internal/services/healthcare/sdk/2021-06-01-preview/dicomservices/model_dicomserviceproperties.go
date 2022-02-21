package dicomservices

type DicomServiceProperties struct {
	AuthenticationConfiguration *DicomServiceAuthenticationConfiguration `json:"authenticationConfiguration,omitempty"`
	ProvisioningState           *ProvisioningState                       `json:"provisioningState,omitempty"`
	ServiceUrl                  *string                                  `json:"serviceUrl,omitempty"`
}
