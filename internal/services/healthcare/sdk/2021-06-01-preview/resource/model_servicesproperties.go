package resource

type ServicesProperties struct {
	AccessPolicies              *[]ServiceAccessPolicyEntry             `json:"accessPolicies,omitempty"`
	AcrConfiguration            *ServiceAcrConfigurationInfo            `json:"acrConfiguration,omitempty"`
	AuthenticationConfiguration *ServiceAuthenticationConfigurationInfo `json:"authenticationConfiguration,omitempty"`
	CorsConfiguration           *ServiceCorsConfigurationInfo           `json:"corsConfiguration,omitempty"`
	CosmosDbConfiguration       *ServiceCosmosDbConfigurationInfo       `json:"cosmosDbConfiguration,omitempty"`
	ExportConfiguration         *ServiceExportConfigurationInfo         `json:"exportConfiguration,omitempty"`
	PrivateEndpointConnections  *[]PrivateEndpointConnection            `json:"privateEndpointConnections,omitempty"`
	ProvisioningState           *ProvisioningState                      `json:"provisioningState,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                    `json:"publicNetworkAccess,omitempty"`
}
