// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_cassandra_table

import (
	"math"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbCassandraTableSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"cassandra_keyspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateCassandraKeyspaceID,
		},

		"default_ttl": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"analytical_storage_ttl": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.IntBetween(-1, math.MaxInt32),
				validation.IntNotInSlice([]int{0}),
			),
		},

		"schema": common.CassandraTableSchemaPropertySchema(),

		"throughput": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.CosmosThroughput,
		},

		"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),
	}
}

func expandTableSchema(d *pluginsdk.ResourceData) *cosmosdb.CassandraSchema {
	i := d.Get("schema").([]any)

	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]any)

	cassandraSchema := cosmosdb.CassandraSchema{}

	if v, ok := input["column"].([]any); ok {
		cassandraSchema.Columns = expandTableSchemaColumns(v)
	}

	if v, ok := input["partition_key"].([]any); ok {
		cassandraSchema.PartitionKeys = expandTableSchemaPartitionKeys(v)
	}

	if v, ok := input["cluster_key"].([]any); ok {
		cassandraSchema.ClusterKeys = expandTableSchemaClusterKeys(v)
	}

	return &cassandraSchema
}

func flattenTableSchema(input *cosmosdb.CassandraSchema) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	result := make(map[string]any)
	result["column"] = flattenTableSchemaColumns(input.Columns)
	result["partition_key"] = flattenTableSchemaPartitionKeys(input.PartitionKeys)
	result["cluster_key"] = flattenTableSchemaClusterKeys(input.ClusterKeys)

	results = append(results, result)
	return results
}
