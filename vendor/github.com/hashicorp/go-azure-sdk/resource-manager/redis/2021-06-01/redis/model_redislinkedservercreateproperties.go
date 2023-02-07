package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisLinkedServerCreateProperties struct {
	LinkedRedisCacheId       string          `json:"linkedRedisCacheId"`
	LinkedRedisCacheLocation string          `json:"linkedRedisCacheLocation"`
	ServerRole               ReplicationRole `json:"serverRole"`
}
