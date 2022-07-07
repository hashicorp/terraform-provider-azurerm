package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProperties struct {
	Containers               []Container                           `json:"containers"`
	Diagnostics              *ContainerGroupDiagnostics            `json:"diagnostics,omitempty"`
	DnsConfig                *DnsConfiguration                     `json:"dnsConfig,omitempty"`
	EncryptionProperties     *EncryptionProperties                 `json:"encryptionProperties,omitempty"`
	ImageRegistryCredentials *[]ImageRegistryCredential            `json:"imageRegistryCredentials,omitempty"`
	InitContainers           *[]InitContainerDefinition            `json:"initContainers,omitempty"`
	InstanceView             *ContainerGroupPropertiesInstanceView `json:"instanceView,omitempty"`
	IpAddress                *IpAddress                            `json:"ipAddress,omitempty"`
	NetworkProfile           *ContainerGroupNetworkProfile         `json:"networkProfile,omitempty"`
	OsType                   OperatingSystemTypes                  `json:"osType"`
	ProvisioningState        *string                               `json:"provisioningState,omitempty"`
	RestartPolicy            *ContainerGroupRestartPolicy          `json:"restartPolicy,omitempty"`
	Sku                      *ContainerGroupSku                    `json:"sku,omitempty"`
	Volumes                  *[]Volume                             `json:"volumes,omitempty"`
}
