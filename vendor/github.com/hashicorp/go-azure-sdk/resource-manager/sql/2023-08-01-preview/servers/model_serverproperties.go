package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerProperties struct {
	AdministratorLogin            *string                            `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword    *string                            `json:"administratorLoginPassword,omitempty"`
	Administrators                *ServerExternalAdministrator       `json:"administrators,omitempty"`
	ExternalGovernanceStatus      *ExternalGovernanceStatus          `json:"externalGovernanceStatus,omitempty"`
	FederatedClientId             *string                            `json:"federatedClientId,omitempty"`
	FullyQualifiedDomainName      *string                            `json:"fullyQualifiedDomainName,omitempty"`
	IsIPv6Enabled                 *ServerNetworkAccessFlag           `json:"isIPv6Enabled,omitempty"`
	KeyId                         *string                            `json:"keyId,omitempty"`
	MinimalTlsVersion             *MinimalTlsVersion                 `json:"minimalTlsVersion,omitempty"`
	PrimaryUserAssignedIdentityId *string                            `json:"primaryUserAssignedIdentityId,omitempty"`
	PrivateEndpointConnections    *[]ServerPrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	PublicNetworkAccess           *ServerPublicNetworkAccessFlag     `json:"publicNetworkAccess,omitempty"`
	RestrictOutboundNetworkAccess *ServerNetworkAccessFlag           `json:"restrictOutboundNetworkAccess,omitempty"`
	State                         *string                            `json:"state,omitempty"`
	Version                       *string                            `json:"version,omitempty"`
	WorkspaceFeature              *ServerWorkspaceFeature            `json:"workspaceFeature,omitempty"`
}
