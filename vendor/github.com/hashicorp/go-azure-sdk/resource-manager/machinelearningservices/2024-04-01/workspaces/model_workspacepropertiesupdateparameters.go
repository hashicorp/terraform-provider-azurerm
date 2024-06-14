package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePropertiesUpdateParameters struct {
	ApplicationInsights             *string                          `json:"applicationInsights,omitempty"`
	ContainerRegistry               *string                          `json:"containerRegistry,omitempty"`
	Description                     *string                          `json:"description,omitempty"`
	EnableDataIsolation             *bool                            `json:"enableDataIsolation,omitempty"`
	Encryption                      *EncryptionUpdateProperties      `json:"encryption,omitempty"`
	FeatureStoreSettings            *FeatureStoreSettings            `json:"featureStoreSettings,omitempty"`
	FriendlyName                    *string                          `json:"friendlyName,omitempty"`
	ImageBuildCompute               *string                          `json:"imageBuildCompute,omitempty"`
	ManagedNetwork                  *ManagedNetworkSettings          `json:"managedNetwork,omitempty"`
	PrimaryUserAssignedIdentity     *string                          `json:"primaryUserAssignedIdentity,omitempty"`
	PublicNetworkAccess             *PublicNetworkAccess             `json:"publicNetworkAccess,omitempty"`
	ServerlessComputeSettings       *ServerlessComputeSettings       `json:"serverlessComputeSettings,omitempty"`
	ServiceManagedResourcesSettings *ServiceManagedResourcesSettings `json:"serviceManagedResourcesSettings,omitempty"`
	V1LegacyMode                    *bool                            `json:"v1LegacyMode,omitempty"`
}
