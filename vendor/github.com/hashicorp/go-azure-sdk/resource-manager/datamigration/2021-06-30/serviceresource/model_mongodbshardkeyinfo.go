package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbShardKeyInfo struct {
	Fields   []MongoDbShardKeyField `json:"fields"`
	IsUnique bool                   `json:"isUnique"`
}
