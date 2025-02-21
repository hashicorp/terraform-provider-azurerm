package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	AgentsEndpointUri               *string                          `json:"agentsEndpointUri,omitempty"`
	AllowPublicAccessWhenBehindVnet *bool                            `json:"allowPublicAccessWhenBehindVnet,omitempty"`
	AllowRoleAssignmentOnRG         *bool                            `json:"allowRoleAssignmentOnRG,omitempty"`
	ApplicationInsights             *string                          `json:"applicationInsights,omitempty"`
	AssociatedWorkspaces            *[]string                        `json:"associatedWorkspaces,omitempty"`
	ContainerRegistries             *[]string                        `json:"containerRegistries,omitempty"`
	ContainerRegistry               *string                          `json:"containerRegistry,omitempty"`
	Description                     *string                          `json:"description,omitempty"`
	DiscoveryURL                    *string                          `json:"discoveryUrl,omitempty"`
	EnableDataIsolation             *bool                            `json:"enableDataIsolation,omitempty"`
	EnableServiceSideCMKEncryption  *bool                            `json:"enableServiceSideCMKEncryption,omitempty"`
	EnableSimplifiedCmk             *bool                            `json:"enableSimplifiedCmk,omitempty"`
	EnableSoftwareBillOfMaterials   *bool                            `json:"enableSoftwareBillOfMaterials,omitempty"`
	Encryption                      *EncryptionProperty              `json:"encryption,omitempty"`
	ExistingWorkspaces              *[]string                        `json:"existingWorkspaces,omitempty"`
	FeatureStoreSettings            *FeatureStoreSettings            `json:"featureStoreSettings,omitempty"`
	FriendlyName                    *string                          `json:"friendlyName,omitempty"`
	HbiWorkspace                    *bool                            `json:"hbiWorkspace,omitempty"`
	HubResourceId                   *string                          `json:"hubResourceId,omitempty"`
	IPAllowlist                     *[]string                        `json:"ipAllowlist,omitempty"`
	ImageBuildCompute               *string                          `json:"imageBuildCompute,omitempty"`
	KeyVault                        *string                          `json:"keyVault,omitempty"`
	KeyVaults                       *[]string                        `json:"keyVaults,omitempty"`
	ManagedNetwork                  *ManagedNetworkSettings          `json:"managedNetwork,omitempty"`
	MlFlowTrackingUri               *string                          `json:"mlFlowTrackingUri,omitempty"`
	NetworkAcls                     *NetworkAcls                     `json:"networkAcls,omitempty"`
	NotebookInfo                    *NotebookResourceInfo            `json:"notebookInfo,omitempty"`
	PrimaryUserAssignedIdentity     *string                          `json:"primaryUserAssignedIdentity,omitempty"`
	PrivateEndpointConnections      *[]PrivateEndpointConnection     `json:"privateEndpointConnections,omitempty"`
	PrivateLinkCount                *int64                           `json:"privateLinkCount,omitempty"`
	ProvisionNetworkNow             *bool                            `json:"provisionNetworkNow,omitempty"`
	ProvisioningState               *ProvisioningState               `json:"provisioningState,omitempty"`
	PublicNetworkAccess             *PublicNetworkAccessType         `json:"publicNetworkAccess,omitempty"`
	ServerlessComputeSettings       *ServerlessComputeSettings       `json:"serverlessComputeSettings,omitempty"`
	ServiceManagedResourcesSettings *ServiceManagedResourcesSettings `json:"serviceManagedResourcesSettings,omitempty"`
	ServiceProvisionedResourceGroup *string                          `json:"serviceProvisionedResourceGroup,omitempty"`
	SharedPrivateLinkResources      *[]SharedPrivateLinkResource     `json:"sharedPrivateLinkResources,omitempty"`
	SoftDeleteRetentionInDays       *int64                           `json:"softDeleteRetentionInDays,omitempty"`
	StorageAccount                  *string                          `json:"storageAccount,omitempty"`
	StorageAccounts                 *[]string                        `json:"storageAccounts,omitempty"`
	StorageHnsEnabled               *bool                            `json:"storageHnsEnabled,omitempty"`
	SystemDatastoresAuthMode        *SystemDatastoresAuthMode        `json:"systemDatastoresAuthMode,omitempty"`
	TenantId                        *string                          `json:"tenantId,omitempty"`
	V1LegacyMode                    *bool                            `json:"v1LegacyMode,omitempty"`
	WorkspaceHubConfig              *WorkspaceHubConfig              `json:"workspaceHubConfig,omitempty"`
	WorkspaceId                     *string                          `json:"workspaceId,omitempty"`
}
