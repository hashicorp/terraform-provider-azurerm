package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBCollectionResource struct {
	AnalyticalStorageTtl *int64                 `json:"analyticalStorageTtl,omitempty"`
	CreateMode           *CreateMode            `json:"createMode,omitempty"`
	Id                   string                 `json:"id"`
	Indexes              *[]MongoIndex          `json:"indexes,omitempty"`
	RestoreParameters    *RestoreParametersBase `json:"restoreParameters,omitempty"`
	ShardKey             *map[string]string     `json:"shardKey,omitempty"`
}
