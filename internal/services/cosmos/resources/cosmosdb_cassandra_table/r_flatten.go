// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func flattenTableSchemaColumns(input *[]cosmosdb.Column) []any {
	if input == nil {
		return []any{}
	}

	columns := make([]any, 0)

	for _, v := range *input {
		name := pointer.From(v.Name)
		typeStr := pointer.From(v.Type)
		columns = append(columns, map[string]any{
			"name": name,
			"type": typeStr,
		})
	}

	return columns
}

func flattenTableSchemaPartitionKeys(input *[]cosmosdb.CassandraPartitionKey) []any {
	if input == nil {
		return []any{}
	}

	keys := make([]any, 0)

	for _, v := range *input {
		keys = append(keys, map[string]any{
			"name": pointer.From(v.Name),
		})
	}

	return keys
}

func flattenTableSchemaClusterKeys(input *[]cosmosdb.ClusterKey) []any {
	if input == nil {
		return []any{}
	}

	keys := make([]any, 0)

	for _, v := range *input {
		name := pointer.From(v.Name)
		orderBy := pointer.From(v.OrderBy)
		keys = append(keys, map[string]any{
			"name":     name,
			"order_by": orderBy,
		})
	}

	return keys
}
