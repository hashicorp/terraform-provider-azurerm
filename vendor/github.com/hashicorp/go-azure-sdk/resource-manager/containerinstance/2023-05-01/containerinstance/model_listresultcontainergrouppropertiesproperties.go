package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListResultContainerGroupPropertiesProperties struct {
	ConfidentialComputeProperties *ConfidentialComputeProperties   `json:"confidentialComputeProperties,omitempty"`
	Containers                    []Container                      `json:"containers"`
	Diagnostics                   *ContainerGroupDiagnostics       `json:"diagnostics,omitempty"`
	DnsConfig                     *DnsConfiguration                `json:"dnsConfig,omitempty"`
	EncryptionProperties          *EncryptionProperties            `json:"encryptionProperties,omitempty"`
	Extensions                    *[]DeploymentExtensionSpec       `json:"extensions,omitempty"`
	IPAddress                     *IPAddress                       `json:"ipAddress,omitempty"`
	ImageRegistryCredentials      *[]ImageRegistryCredential       `json:"imageRegistryCredentials,omitempty"`
	InitContainers                *[]InitContainerDefinition       `json:"initContainers,omitempty"`
	OsType                        OperatingSystemTypes             `json:"osType"`
	Priority                      *ContainerGroupPriority          `json:"priority,omitempty"`
	ProvisioningState             *ContainerGroupProvisioningState `json:"provisioningState,omitempty"`
	RestartPolicy                 *ContainerGroupRestartPolicy     `json:"restartPolicy,omitempty"`
	Sku                           *ContainerGroupSku               `json:"sku,omitempty"`
	SubnetIds                     *[]ContainerGroupSubnetId        `json:"subnetIds,omitempty"`
	Volumes                       *[]Volume                        `json:"volumes,omitempty"`
}
