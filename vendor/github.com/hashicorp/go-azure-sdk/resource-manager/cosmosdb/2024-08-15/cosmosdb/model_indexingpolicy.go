package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexingPolicy struct {
	Automatic        *bool              `json:"automatic,omitempty"`
	CompositeIndexes *[][]CompositePath `json:"compositeIndexes,omitempty"`
	ExcludedPaths    *[]ExcludedPath    `json:"excludedPaths,omitempty"`
	IncludedPaths    *[]IncludedPath    `json:"includedPaths,omitempty"`
	IndexingMode     *IndexingMode      `json:"indexingMode,omitempty"`
	SpatialIndexes   *[]SpatialSpec     `json:"spatialIndexes,omitempty"`
}
