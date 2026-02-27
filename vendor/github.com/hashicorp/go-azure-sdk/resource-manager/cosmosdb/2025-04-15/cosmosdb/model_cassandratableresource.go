package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraTableResource struct {
	AnalyticalStorageTtl *int64           `json:"analyticalStorageTtl,omitempty"`
	DefaultTtl           *int64           `json:"defaultTtl,omitempty"`
	Id                   string           `json:"id"`
	Schema               *CassandraSchema `json:"schema,omitempty"`
}
