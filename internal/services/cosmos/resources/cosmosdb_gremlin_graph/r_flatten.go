// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenAzureRmCosmosDBGremlinGraphIndexingPolicy(input *cosmosdb.IndexingPolicy) []any {
	if input == nil {
		return []any{}
	}
	indexPolicy := make(map[string]any)

	indexPolicy["automatic"] = input.Automatic
	indexPolicy["indexing_mode"] = input.IndexingMode
	indexPolicy["included_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input.IncludedPaths))
	indexPolicy["excluded_paths"] = pluginsdk.NewSet(pluginsdk.HashString, flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input.ExcludedPaths))
	indexPolicy["composite_index"] = common.FlattenCosmosDBIndexingPolicyCompositeIndexes(input.CompositeIndexes)
	indexPolicy["spatial_index"] = common.FlattenCosmosDBIndexingPolicySpatialIndexes(input.SpatialIndexes)

	return []any{indexPolicy}
}

func flattenAzureRmCosmosDBGremlinGraphIncludedPaths(input *[]cosmosdb.IncludedPath) []any {
	if input == nil {
		return []any{}
	}

	includedPaths := make([]any, 0)
	for _, includedPath := range *input {
		if includedPath.Path == nil {
			continue
		}

		includedPaths = append(includedPaths, *includedPath.Path)
	}

	return includedPaths
}

func flattenAzureRmCosmosDBGremlinGraphExcludedPaths(input *[]cosmosdb.ExcludedPath) []any {
	if input == nil {
		return []any{}
	}

	excludedPaths := make([]any, 0)
	for _, excludedPath := range *input {
		if excludedPath.Path == nil {
			continue
		}

		excludedPaths = append(excludedPaths, *excludedPath.Path)
	}

	return excludedPaths
}

func flattenCosmosGremlinGraphUniqueKeys(keys *[]cosmosdb.UniqueKey) *[]map[string]any {
	if keys == nil {
		return nil
	}

	slice := make([]map[string]any, 0)
	for _, k := range *keys {
		if k.Paths == nil {
			continue
		}

		slice = append(slice, map[string]any{
			"paths": *k.Paths,
		})
	}

	return &slice
}
