package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	AllowPublicAccessWhenBehindVnet *bool                            `json:"allowPublicAccessWhenBehindVnet,omitempty"`
	ApplicationInsights             *string                          `json:"applicationInsights,omitempty"`
	AssociatedWorkspaces            *[]string                        `json:"associatedWorkspaces,omitempty"`
	ContainerRegistry               *string                          `json:"containerRegistry,omitempty"`
	Description                     *string                          `json:"description,omitempty"`
	DiscoveryURL                    *string                          `json:"discoveryUrl,omitempty"`
	EnableDataIsolation             *bool                            `json:"enableDataIsolation,omitempty"`
	Encryption                      *EncryptionProperty              `json:"encryption,omitempty"`
	FeatureStoreSettings            *FeatureStoreSettings            `json:"featureStoreSettings,omitempty"`
	FriendlyName                    *string                          `json:"friendlyName,omitempty"`
	HbiWorkspace                    *bool                            `json:"hbiWorkspace,omitempty"`
	HubResourceId                   *string                          `json:"hubResourceId,omitempty"`
	ImageBuildCompute               *string                          `json:"imageBuildCompute,omitempty"`
	KeyVault                        *string                          `json:"keyVault,omitempty"`
	ManagedNetwork                  *ManagedNetworkSettings          `json:"managedNetwork,omitempty"`
	MlFlowTrackingUri               *string                          `json:"mlFlowTrackingUri,omitempty"`
	NotebookInfo                    *NotebookResourceInfo            `json:"notebookInfo,omitempty"`
	PrimaryUserAssignedIdentity     *string                          `json:"primaryUserAssignedIdentity,omitempty"`
	PrivateEndpointConnections      *[]PrivateEndpointConnection     `json:"privateEndpointConnections,omitempty"`
	PrivateLinkCount                *int64                           `json:"privateLinkCount,omitempty"`
	ProvisioningState               *ProvisioningState               `json:"provisioningState,omitempty"`
	PublicNetworkAccess             *PublicNetworkAccess             `json:"publicNetworkAccess,omitempty"`
	ServerlessComputeSettings       *ServerlessComputeSettings       `json:"serverlessComputeSettings,omitempty"`
	ServiceManagedResourcesSettings *ServiceManagedResourcesSettings `json:"serviceManagedResourcesSettings,omitempty"`
	ServiceProvisionedResourceGroup *string                          `json:"serviceProvisionedResourceGroup,omitempty"`
	SharedPrivateLinkResources      *[]SharedPrivateLinkResource     `json:"sharedPrivateLinkResources,omitempty"`
	StorageAccount                  *string                          `json:"storageAccount,omitempty"`
	StorageHnsEnabled               *bool                            `json:"storageHnsEnabled,omitempty"`
	TenantId                        *string                          `json:"tenantId,omitempty"`
	V1LegacyMode                    *bool                            `json:"v1LegacyMode,omitempty"`
	WorkspaceHubConfig              *WorkspaceHubConfig              `json:"workspaceHubConfig,omitempty"`
	WorkspaceId                     *string                          `json:"workspaceId,omitempty"`
}
