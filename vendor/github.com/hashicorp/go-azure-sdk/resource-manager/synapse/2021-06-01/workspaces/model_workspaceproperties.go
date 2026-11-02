package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	AdlaResourceId                   *string                           `json:"adlaResourceId,omitempty"`
	AzureADOnlyAuthentication        *bool                             `json:"azureADOnlyAuthentication,omitempty"`
	ConnectivityEndpoints            *map[string]string                `json:"connectivityEndpoints,omitempty"`
	CspWorkspaceAdminProperties      *CspWorkspaceAdminProperties      `json:"cspWorkspaceAdminProperties,omitempty"`
	DefaultDataLakeStorage           *DataLakeStorageAccountDetails    `json:"defaultDataLakeStorage,omitempty"`
	Encryption                       *EncryptionDetails                `json:"encryption,omitempty"`
	ExtraProperties                  *interface{}                      `json:"extraProperties,omitempty"`
	ManagedResourceGroupName         *string                           `json:"managedResourceGroupName,omitempty"`
	ManagedVirtualNetwork            *string                           `json:"managedVirtualNetwork,omitempty"`
	ManagedVirtualNetworkSettings    *ManagedVirtualNetworkSettings    `json:"managedVirtualNetworkSettings,omitempty"`
	PrivateEndpointConnections       *[]PrivateEndpointConnection      `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                *string                           `json:"provisioningState,omitempty"`
	PublicNetworkAccess              *WorkspacePublicNetworkAccess     `json:"publicNetworkAccess,omitempty"`
	PurviewConfiguration             *PurviewConfiguration             `json:"purviewConfiguration,omitempty"`
	Settings                         *map[string]interface{}           `json:"settings,omitempty"`
	SqlAdministratorLogin            *string                           `json:"sqlAdministratorLogin,omitempty"`
	SqlAdministratorLoginPassword    *string                           `json:"sqlAdministratorLoginPassword,omitempty"`
	TrustedServiceBypassEnabled      *bool                             `json:"trustedServiceBypassEnabled,omitempty"`
	VirtualNetworkProfile            *VirtualNetworkProfile            `json:"virtualNetworkProfile,omitempty"`
	WorkspaceRepositoryConfiguration *WorkspaceRepositoryConfiguration `json:"workspaceRepositoryConfiguration,omitempty"`
	WorkspaceUID                     *string                           `json:"workspaceUID,omitempty"`
}
