package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmazonRdsForOraclePartitionSettings struct {
	PartitionColumnName *interface{} `json:"partitionColumnName,omitempty"`
	PartitionLowerBound *interface{} `json:"partitionLowerBound,omitempty"`
	PartitionNames      *interface{} `json:"partitionNames,omitempty"`
	PartitionUpperBound *interface{} `json:"partitionUpperBound,omitempty"`
}
