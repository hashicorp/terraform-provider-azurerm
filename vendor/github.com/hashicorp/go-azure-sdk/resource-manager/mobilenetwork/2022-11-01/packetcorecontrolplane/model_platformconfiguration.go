package packetcorecontrolplane

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlatformConfiguration struct {
	AzureStackEdgeDevice  *AzureStackEdgeDeviceResourceId   `json:"azureStackEdgeDevice,omitempty"`
	AzureStackEdgeDevices *[]AzureStackEdgeDeviceResourceId `json:"azureStackEdgeDevices,omitempty"`
	AzureStackHciCluster  *AzureStackHCIClusterResourceId   `json:"azureStackHciCluster,omitempty"`
	ConnectedCluster      *ConnectedClusterResourceId       `json:"connectedCluster,omitempty"`
	CustomLocation        *CustomLocationResourceId         `json:"customLocation,omitempty"`
	Type                  PlatformType                      `json:"type"`
}
