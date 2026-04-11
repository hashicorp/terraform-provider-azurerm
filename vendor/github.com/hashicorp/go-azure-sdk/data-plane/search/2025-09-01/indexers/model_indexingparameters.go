package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexingParameters struct {
	BatchSize              *int64                           `json:"batchSize,omitempty"`
	Configuration          *IndexingParametersConfiguration `json:"configuration,omitempty"`
	MaxFailedItems         *int64                           `json:"maxFailedItems,omitempty"`
	MaxFailedItemsPerBatch *int64                           `json:"maxFailedItemsPerBatch,omitempty"`
}
