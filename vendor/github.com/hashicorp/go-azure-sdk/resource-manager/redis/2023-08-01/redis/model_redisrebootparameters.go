package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisRebootParameters struct {
	Ports      *[]int64    `json:"ports,omitempty"`
	RebootType *RebootType `json:"rebootType,omitempty"`
	ShardId    *int64      `json:"shardId,omitempty"`
}
