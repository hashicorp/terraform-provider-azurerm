package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCreateProperties struct {
	DisableAccessKeyAuthentication *bool                                    `json:"disableAccessKeyAuthentication,omitempty"`
	EnableNonSslPort               *bool                                    `json:"enableNonSslPort,omitempty"`
	MinimumTlsVersion              *TlsVersion                              `json:"minimumTlsVersion,omitempty"`
	PublicNetworkAccess            *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
	RedisConfiguration             *RedisCommonPropertiesRedisConfiguration `json:"redisConfiguration,omitempty"`
	RedisVersion                   *string                                  `json:"redisVersion,omitempty"`
	ReplicasPerMaster              *int64                                   `json:"replicasPerMaster,omitempty"`
	ReplicasPerPrimary             *int64                                   `json:"replicasPerPrimary,omitempty"`
	ShardCount                     *int64                                   `json:"shardCount,omitempty"`
	Sku                            Sku                                      `json:"sku"`
	StaticIP                       *string                                  `json:"staticIP,omitempty"`
	SubnetId                       *string                                  `json:"subnetId,omitempty"`
	TenantSettings                 *map[string]string                       `json:"tenantSettings,omitempty"`
	UpdateChannel                  *UpdateChannel                           `json:"updateChannel,omitempty"`
}
