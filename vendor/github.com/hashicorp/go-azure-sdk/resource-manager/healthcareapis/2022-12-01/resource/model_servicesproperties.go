package resource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesProperties struct {
	AccessPolicies              *[]ServiceAccessPolicyEntry             `json:"accessPolicies,omitempty"`
	AcrConfiguration            *ServiceAcrConfigurationInfo            `json:"acrConfiguration,omitempty"`
	AuthenticationConfiguration *ServiceAuthenticationConfigurationInfo `json:"authenticationConfiguration,omitempty"`
	CorsConfiguration           *ServiceCorsConfigurationInfo           `json:"corsConfiguration,omitempty"`
	CosmosDbConfiguration       *ServiceCosmosDbConfigurationInfo       `json:"cosmosDbConfiguration,omitempty"`
	ExportConfiguration         *ServiceExportConfigurationInfo         `json:"exportConfiguration,omitempty"`
	ImportConfiguration         *ServiceImportConfigurationInfo         `json:"importConfiguration,omitempty"`
	PrivateEndpointConnections  *[]PrivateEndpointConnection            `json:"privateEndpointConnections,omitempty"`
	ProvisioningState           *ProvisioningState                      `json:"provisioningState,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                    `json:"publicNetworkAccess,omitempty"`
}
