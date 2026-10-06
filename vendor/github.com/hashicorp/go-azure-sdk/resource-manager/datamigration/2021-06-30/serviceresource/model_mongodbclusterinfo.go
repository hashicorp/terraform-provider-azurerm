package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbClusterInfo struct {
	Databases        []MongoDbDatabaseInfo `json:"databases"`
	SupportsSharding bool                  `json:"supportsSharding"`
	Type             MongoDbClusterType    `json:"type"`
	Version          string                `json:"version"`
}
