package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoIndexOptions struct {
	ExpireAfterSeconds *int64 `json:"expireAfterSeconds,omitempty"`
	Unique             *bool  `json:"unique,omitempty"`
}
