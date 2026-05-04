package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraSchema struct {
	ClusterKeys   *[]ClusterKey            `json:"clusterKeys,omitempty"`
	Columns       *[]Column                `json:"columns,omitempty"`
	PartitionKeys *[]CassandraPartitionKey `json:"partitionKeys,omitempty"`
}
