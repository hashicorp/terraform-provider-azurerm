package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbCollectionInfo struct {
	AverageDocumentSize int64                `json:"averageDocumentSize"`
	DataSize            int64                `json:"dataSize"`
	DatabaseName        string               `json:"databaseName"`
	DocumentCount       int64                `json:"documentCount"`
	IsCapped            bool                 `json:"isCapped"`
	IsSystemCollection  bool                 `json:"isSystemCollection"`
	IsView              bool                 `json:"isView"`
	Name                string               `json:"name"`
	QualifiedName       string               `json:"qualifiedName"`
	ShardKey            *MongoDbShardKeyInfo `json:"shardKey,omitempty"`
	SupportsSharding    bool                 `json:"supportsSharding"`
	ViewOf              *string              `json:"viewOf,omitempty"`
}
