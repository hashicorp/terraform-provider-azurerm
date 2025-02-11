package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSite struct {
	AllowConfigFileUpdates      *bool                                                     `json:"allowConfigFileUpdates,omitempty"`
	Branch                      *string                                                   `json:"branch,omitempty"`
	BuildProperties             *StaticSiteBuildProperties                                `json:"buildProperties,omitempty"`
	ContentDistributionEndpoint *string                                                   `json:"contentDistributionEndpoint,omitempty"`
	CustomDomains               *[]string                                                 `json:"customDomains,omitempty"`
	DatabaseConnections         *[]DatabaseConnectionOverview                             `json:"databaseConnections,omitempty"`
	DefaultHostname             *string                                                   `json:"defaultHostname,omitempty"`
	EnterpriseGradeCdnStatus    *EnterpriseGradeCdnStatus                                 `json:"enterpriseGradeCdnStatus,omitempty"`
	KeyVaultReferenceIdentity   *string                                                   `json:"keyVaultReferenceIdentity,omitempty"`
	LinkedBackends              *[]StaticSiteLinkedBackend                                `json:"linkedBackends,omitempty"`
	PrivateEndpointConnections  *[]ResponseMessageEnvelopeRemotePrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	Provider                    *string                                                   `json:"provider,omitempty"`
	PublicNetworkAccess         *string                                                   `json:"publicNetworkAccess,omitempty"`
	RepositoryToken             *string                                                   `json:"repositoryToken,omitempty"`
	RepositoryURL               *string                                                   `json:"repositoryUrl,omitempty"`
	StagingEnvironmentPolicy    *StagingEnvironmentPolicy                                 `json:"stagingEnvironmentPolicy,omitempty"`
	TemplateProperties          *StaticSiteTemplateOptions                                `json:"templateProperties,omitempty"`
	UserProvidedFunctionApps    *[]StaticSiteUserProvidedFunctionApp                      `json:"userProvidedFunctionApps,omitempty"`
}
