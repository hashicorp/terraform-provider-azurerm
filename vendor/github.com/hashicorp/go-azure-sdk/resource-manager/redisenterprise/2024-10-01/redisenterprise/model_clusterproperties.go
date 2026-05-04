package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	Encryption                 *ClusterPropertiesEncryption `json:"encryption,omitempty"`
	HostName                   *string                      `json:"hostName,omitempty"`
	MinimumTlsVersion          *TlsVersion                  `json:"minimumTlsVersion,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	RedisVersion               *string                      `json:"redisVersion,omitempty"`
	ResourceState              *ResourceState               `json:"resourceState,omitempty"`
}
