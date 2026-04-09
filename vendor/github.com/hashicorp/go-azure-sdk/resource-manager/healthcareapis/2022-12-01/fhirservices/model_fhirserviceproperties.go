package fhirservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FhirServiceProperties struct {
	AccessPolicies                     *[]FhirServiceAccessPolicyEntry         `json:"accessPolicies,omitempty"`
	AcrConfiguration                   *FhirServiceAcrConfiguration            `json:"acrConfiguration,omitempty"`
	AuthenticationConfiguration        *FhirServiceAuthenticationConfiguration `json:"authenticationConfiguration,omitempty"`
	CorsConfiguration                  *FhirServiceCorsConfiguration           `json:"corsConfiguration,omitempty"`
	EventState                         *ServiceEventState                      `json:"eventState,omitempty"`
	ExportConfiguration                *FhirServiceExportConfiguration         `json:"exportConfiguration,omitempty"`
	ImplementationGuidesConfiguration  *ImplementationGuidesConfiguration      `json:"implementationGuidesConfiguration,omitempty"`
	ImportConfiguration                *FhirServiceImportConfiguration         `json:"importConfiguration,omitempty"`
	PrivateEndpointConnections         *[]PrivateEndpointConnection            `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                  *ProvisioningState                      `json:"provisioningState,omitempty"`
	PublicNetworkAccess                *PublicNetworkAccess                    `json:"publicNetworkAccess,omitempty"`
	ResourceVersionPolicyConfiguration *ResourceVersionPolicyConfiguration     `json:"resourceVersionPolicyConfiguration,omitempty"`
}
