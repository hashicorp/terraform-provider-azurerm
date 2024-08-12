package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisProperties struct {
	AccessKeys                     *RedisAccessKeys                         `json:"accessKeys,omitempty"`
	DisableAccessKeyAuthentication *bool                                    `json:"disableAccessKeyAuthentication,omitempty"`
	EnableNonSslPort               *bool                                    `json:"enableNonSslPort,omitempty"`
	HostName                       *string                                  `json:"hostName,omitempty"`
	Instances                      *[]RedisInstanceDetails                  `json:"instances,omitempty"`
	LinkedServers                  *[]RedisLinkedServer                     `json:"linkedServers,omitempty"`
	MinimumTlsVersion              *TlsVersion                              `json:"minimumTlsVersion,omitempty"`
	Port                           *int64                                   `json:"port,omitempty"`
	PrivateEndpointConnections     *[]PrivateEndpointConnection             `json:"privateEndpointConnections,omitempty"`
	ProvisioningState              *ProvisioningState                       `json:"provisioningState,omitempty"`
	PublicNetworkAccess            *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
	RedisConfiguration             *RedisCommonPropertiesRedisConfiguration `json:"redisConfiguration,omitempty"`
	RedisVersion                   *string                                  `json:"redisVersion,omitempty"`
	ReplicasPerMaster              *int64                                   `json:"replicasPerMaster,omitempty"`
	ReplicasPerPrimary             *int64                                   `json:"replicasPerPrimary,omitempty"`
	ShardCount                     *int64                                   `json:"shardCount,omitempty"`
	Sku                            Sku                                      `json:"sku"`
	SslPort                        *int64                                   `json:"sslPort,omitempty"`
	StaticIP                       *string                                  `json:"staticIP,omitempty"`
	SubnetId                       *string                                  `json:"subnetId,omitempty"`
	TenantSettings                 *map[string]string                       `json:"tenantSettings,omitempty"`
	UpdateChannel                  *UpdateChannel                           `json:"updateChannel,omitempty"`
}
