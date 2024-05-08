package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolResourceProperties struct {
	AksClusterProfile           *AksClusterProfile              `json:"aksClusterProfile,omitempty"`
	AksManagedResourceGroupName *string                         `json:"aksManagedResourceGroupName,omitempty"`
	ClusterPoolProfile          *ClusterPoolProfile             `json:"clusterPoolProfile,omitempty"`
	ComputeProfile              ClusterPoolComputeProfile       `json:"computeProfile"`
	DeploymentId                *string                         `json:"deploymentId,omitempty"`
	LogAnalyticsProfile         *ClusterPoolLogAnalyticsProfile `json:"logAnalyticsProfile,omitempty"`
	ManagedResourceGroupName    *string                         `json:"managedResourceGroupName,omitempty"`
	NetworkProfile              *ClusterPoolNetworkProfile      `json:"networkProfile,omitempty"`
	ProvisioningState           *ProvisioningStatus             `json:"provisioningState,omitempty"`
	Status                      *string                         `json:"status,omitempty"`
}
