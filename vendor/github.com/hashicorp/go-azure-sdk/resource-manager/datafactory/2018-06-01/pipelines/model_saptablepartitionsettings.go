package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapTablePartitionSettings struct {
	MaxPartitionsNumber *int64  `json:"maxPartitionsNumber,omitempty"`
	PartitionColumnName *string `json:"partitionColumnName,omitempty"`
	PartitionLowerBound *string `json:"partitionLowerBound,omitempty"`
	PartitionUpperBound *string `json:"partitionUpperBound,omitempty"`
}
