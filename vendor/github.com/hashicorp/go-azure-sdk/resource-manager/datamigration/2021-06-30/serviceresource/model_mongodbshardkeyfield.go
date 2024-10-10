package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbShardKeyField struct {
	Name  string               `json:"name"`
	Order MongoDbShardKeyOrder `json:"order"`
}
