package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePropertiesUpdateParameters struct {
	AllowRoleAssignmentOnRG         *bool                            `json:"allowRoleAssignmentOnRG,omitempty"`
	ApplicationInsights             *string                          `json:"applicationInsights,omitempty"`
	ContainerRegistry               *string                          `json:"containerRegistry,omitempty"`
	Description                     *string                          `json:"description,omitempty"`
	EnableDataIsolation             *bool                            `json:"enableDataIsolation,omitempty"`
	EnableSoftwareBillOfMaterials   *bool                            `json:"enableSoftwareBillOfMaterials,omitempty"`
	Encryption                      *EncryptionUpdateProperties      `json:"encryption,omitempty"`
	FeatureStoreSettings            *FeatureStoreSettings            `json:"featureStoreSettings,omitempty"`
	FriendlyName                    *string                          `json:"friendlyName,omitempty"`
	IPAllowlist                     *[]string                        `json:"ipAllowlist,omitempty"`
	ImageBuildCompute               *string                          `json:"imageBuildCompute,omitempty"`
	ManagedNetwork                  *ManagedNetworkSettings          `json:"managedNetwork,omitempty"`
	NetworkAcls                     *NetworkAcls                     `json:"networkAcls,omitempty"`
	PrimaryUserAssignedIdentity     *string                          `json:"primaryUserAssignedIdentity,omitempty"`
	PublicNetworkAccess             *PublicNetworkAccessType         `json:"publicNetworkAccess,omitempty"`
	ServerlessComputeSettings       *ServerlessComputeSettings       `json:"serverlessComputeSettings,omitempty"`
	ServiceManagedResourcesSettings *ServiceManagedResourcesSettings `json:"serviceManagedResourcesSettings,omitempty"`
	SoftDeleteRetentionInDays       *int64                           `json:"softDeleteRetentionInDays,omitempty"`
	SystemDatastoresAuthMode        *SystemDatastoresAuthMode        `json:"systemDatastoresAuthMode,omitempty"`
	V1LegacyMode                    *bool                            `json:"v1LegacyMode,omitempty"`
}
