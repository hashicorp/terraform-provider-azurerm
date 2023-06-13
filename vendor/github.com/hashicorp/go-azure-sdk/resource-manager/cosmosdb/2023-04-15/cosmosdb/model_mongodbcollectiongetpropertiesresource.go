package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBCollectionGetPropertiesResource struct {
	AnalyticalStorageTtl *int64             `json:"analyticalStorageTtl,omitempty"`
	Etag                 *string            `json:"_etag,omitempty"`
	Id                   *string            `json:"id,omitempty"`
	Indexes              *[]MongoIndex      `json:"indexes,omitempty"`
	Rid                  *string            `json:"_rid,omitempty"`
	ShardKey             *map[string]string `json:"shardKey,omitempty"`
	Ts                   *float64           `json:"_ts,omitempty"`
}
