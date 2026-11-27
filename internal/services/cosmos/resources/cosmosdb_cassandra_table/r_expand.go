// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func expandTableSchemaColumns(input []any) *[]cosmosdb.Column {
	columns := make([]cosmosdb.Column, 0)
	for _, col := range input {
		data := col.(map[string]any)
		column := cosmosdb.Column{
			Name: pointer.To(data["name"].(string)),
			Type: pointer.To(data["type"].(string)),
		}
		columns = append(columns, column)
	}

	return &columns
}

func expandTableSchemaPartitionKeys(input []any) *[]cosmosdb.CassandraPartitionKey {
	keys := make([]cosmosdb.CassandraPartitionKey, 0)
	for _, key := range input {
		data := key.(map[string]any)
		k := cosmosdb.CassandraPartitionKey{
			Name: pointer.To(data["name"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}

func expandTableSchemaClusterKeys(input []any) *[]cosmosdb.ClusterKey {
	keys := make([]cosmosdb.ClusterKey, 0)
	for _, key := range input {
		data := key.(map[string]any)
		k := cosmosdb.ClusterKey{
			Name:    pointer.To(data["name"].(string)),
			OrderBy: pointer.To(data["order_by"].(string)),
		}
		keys = append(keys, k)
	}

	return &keys
}
