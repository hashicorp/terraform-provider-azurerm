package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisLinkedServerProperties struct {
	GeoReplicatedPrimaryHostName *string         `json:"geoReplicatedPrimaryHostName,omitempty"`
	LinkedRedisCacheId           string          `json:"linkedRedisCacheId"`
	LinkedRedisCacheLocation     string          `json:"linkedRedisCacheLocation"`
	PrimaryHostName              *string         `json:"primaryHostName,omitempty"`
	ProvisioningState            *string         `json:"provisioningState,omitempty"`
	ServerRole                   ReplicationRole `json:"serverRole"`
}
