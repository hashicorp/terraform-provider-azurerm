package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupPropertiesProperties struct {
	Containers               []Container                                     `json:"containers"`
	Diagnostics              *ContainerGroupDiagnostics                      `json:"diagnostics"`
	DnsConfig                *DnsConfiguration                               `json:"dnsConfig"`
	EncryptionProperties     *EncryptionProperties                           `json:"encryptionProperties"`
	IPAddress                *IPAddress                                      `json:"ipAddress"`
	ImageRegistryCredentials *[]ImageRegistryCredential                      `json:"imageRegistryCredentials,omitempty"`
	InitContainers           *[]InitContainerDefinition                      `json:"initContainers,omitempty"`
	InstanceView             *ContainerGroupPropertiesPropertiesInstanceView `json:"instanceView"`
	OsType                   OperatingSystemTypes                            `json:"osType"`
	ProvisioningState        *string                                         `json:"provisioningState,omitempty"`
	RestartPolicy            *ContainerGroupRestartPolicy                    `json:"restartPolicy,omitempty"`
	Sku                      *ContainerGroupSku                              `json:"sku,omitempty"`
	SubnetIds                *[]ContainerGroupSubnetId                       `json:"subnetIds,omitempty"`
	Volumes                  *[]Volume                                       `json:"volumes,omitempty"`
}
