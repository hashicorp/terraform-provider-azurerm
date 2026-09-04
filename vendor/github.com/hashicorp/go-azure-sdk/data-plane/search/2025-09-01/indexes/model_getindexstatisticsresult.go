package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetIndexStatisticsResult struct {
	DocumentCount   int64 `json:"documentCount"`
	StorageSize     int64 `json:"storageSize"`
	VectorIndexSize int64 `json:"vectorIndexSize"`
}
