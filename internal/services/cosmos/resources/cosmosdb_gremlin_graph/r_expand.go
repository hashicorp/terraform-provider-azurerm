// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAzureRmCosmosDbGremlinGraphIndexingPolicy(d *pluginsdk.ResourceData) *cosmosdb.IndexingPolicy {
	i := d.Get("index_policy").([]any)
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	input := i[0].(map[string]any)
	indexingPolicy := cosmosdb.IndexingMode(strings.ToLower(input["indexing_mode"].(string)))
	policy := &cosmosdb.IndexingPolicy{
		IndexingMode:  &indexingPolicy,
		IncludedPaths: expandAzureRmCosmosDbGremlinGraphIncludedPath(input),
		ExcludedPaths: expandAzureRmCosmosDbGremlinGraphExcludedPath(input),
	}
	if v, ok := input["composite_index"].([]any); ok {
		policy.CompositeIndexes = common.ExpandAzureRmCosmosDBIndexingPolicyCompositeIndexes(v)
	}

	policy.SpatialIndexes = common.ExpandAzureRmCosmosDBIndexingPolicySpatialIndexes(input["spatial_index"].([]any))

	if automatic, ok := input["automatic"].(bool); ok {
		policy.Automatic = pointer.To(automatic)
	}

	return policy
}

func expandAzureRmCosmosDbGremlinGraphIncludedPath(input map[string]any) *[]cosmosdb.IncludedPath {
	includedPath := input["included_paths"].(*pluginsdk.Set).List()
	paths := make([]cosmosdb.IncludedPath, len(includedPath))

	for i, pathConfig := range includedPath {
		attrs := pathConfig.(string)
		paths[i] = cosmosdb.IncludedPath{
			Path: pointer.To(attrs),
		}
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphExcludedPath(input map[string]any) *[]cosmosdb.ExcludedPath {
	excludedPath := input["excluded_paths"].(*pluginsdk.Set).List()
	paths := make([]cosmosdb.ExcludedPath, len(excludedPath))

	for i, pathConfig := range excludedPath {
		attrs := pathConfig.(string)
		paths[i] = cosmosdb.ExcludedPath{
			Path: pointer.To(attrs),
		}
	}

	return &paths
}

func expandAzureRmCosmosDbGremlinGraphUniqueKeys(s *pluginsdk.Set) *[]cosmosdb.UniqueKey {
	i := s.List()
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	keys := make([]cosmosdb.UniqueKey, 0)
	for _, k := range i {
		key := k.(map[string]any)

		paths := key["paths"].(*pluginsdk.Set).List()
		if len(paths) == 0 {
			continue
		}

		keys = append(keys, cosmosdb.UniqueKey{
			Paths: pluginsdk.ExpandStringSlice(paths),
		})
	}

	return &keys
}
