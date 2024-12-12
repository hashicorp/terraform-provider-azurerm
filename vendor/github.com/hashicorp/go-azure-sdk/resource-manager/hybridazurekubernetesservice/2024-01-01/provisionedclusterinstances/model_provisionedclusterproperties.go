package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisionedClusterProperties struct {
	AgentPoolProfiles      *[]NamedAgentPoolProfile                       `json:"agentPoolProfiles,omitempty"`
	AutoScalerProfile      *ProvisionedClusterPropertiesAutoScalerProfile `json:"autoScalerProfile,omitempty"`
	CloudProviderProfile   *CloudProviderProfile                          `json:"cloudProviderProfile,omitempty"`
	ClusterVMAccessProfile *ClusterVMAccessProfile                        `json:"clusterVMAccessProfile,omitempty"`
	ControlPlane           *ControlPlaneProfile                           `json:"controlPlane,omitempty"`
	KubernetesVersion      *string                                        `json:"kubernetesVersion,omitempty"`
	LicenseProfile         *ProvisionedClusterLicenseProfile              `json:"licenseProfile,omitempty"`
	LinuxProfile           *LinuxProfileProperties                        `json:"linuxProfile,omitempty"`
	NetworkProfile         *NetworkProfile                                `json:"networkProfile,omitempty"`
	ProvisioningState      *ResourceProvisioningState                     `json:"provisioningState,omitempty"`
	Status                 *ProvisionedClusterPropertiesStatus            `json:"status,omitempty"`
	StorageProfile         *StorageProfile                                `json:"storageProfile,omitempty"`
}
