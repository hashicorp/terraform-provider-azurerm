package afdcustomdomains

type AFDDomainUpdatePropertiesParameters struct {
	AzureDnsZone                       *ResourceReference        `json:"azureDnsZone,omitempty"`
	PreValidatedCustomDomainResourceId *ResourceReference        `json:"preValidatedCustomDomainResourceId,omitempty"`
	ProfileName                        *string                   `json:"profileName,omitempty"`
	TlsSettings                        *AFDDomainHttpsParameters `json:"tlsSettings,omitempty"`
}
