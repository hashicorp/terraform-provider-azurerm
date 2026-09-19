package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorIndex struct {
	IndexingSearchListSize *int64          `json:"indexingSearchListSize,omitempty"`
	Path                   string          `json:"path"`
	QuantizationByteSize   *int64          `json:"quantizationByteSize,omitempty"`
	Type                   VectorIndexType `json:"type"`
	VectorIndexShardKey    *[]string       `json:"vectorIndexShardKey,omitempty"`
}
