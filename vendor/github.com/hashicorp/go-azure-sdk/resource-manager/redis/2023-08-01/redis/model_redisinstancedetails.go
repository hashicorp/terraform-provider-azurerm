package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisInstanceDetails struct {
	IsMaster   *bool   `json:"isMaster,omitempty"`
	IsPrimary  *bool   `json:"isPrimary,omitempty"`
	NonSslPort *int64  `json:"nonSslPort,omitempty"`
	ShardId    *int64  `json:"shardId,omitempty"`
	SslPort    *int64  `json:"sslPort,omitempty"`
	Zone       *string `json:"zone,omitempty"`
}
