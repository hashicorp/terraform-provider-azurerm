package afdcustomdomains

type AFDDomainProperties struct {
	AzureDnsZone                       *ResourceReference          `json:"azureDnsZone,omitempty"`
	DeploymentStatus                   *DeploymentStatus           `json:"deploymentStatus,omitempty"`
	DomainValidationState              *DomainValidationState      `json:"domainValidationState,omitempty"`
	HostName                           string                      `json:"hostName"`
	PreValidatedCustomDomainResourceId *ResourceReference          `json:"preValidatedCustomDomainResourceId,omitempty"`
	ProfileName                        *string                     `json:"profileName,omitempty"`
	ProvisioningState                  *AfdProvisioningState       `json:"provisioningState,omitempty"`
	TlsSettings                        *AFDDomainHttpsParameters   `json:"tlsSettings,omitempty"`
	ValidationProperties               *DomainValidationProperties `json:"validationProperties,omitempty"`
}
