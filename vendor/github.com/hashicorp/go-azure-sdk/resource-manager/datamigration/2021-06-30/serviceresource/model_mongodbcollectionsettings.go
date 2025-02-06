package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbCollectionSettings struct {
	CanDelete *bool                   `json:"canDelete,omitempty"`
	ShardKey  *MongoDbShardKeySetting `json:"shardKey,omitempty"`
	TargetRUs *int64                  `json:"targetRUs,omitempty"`
}
