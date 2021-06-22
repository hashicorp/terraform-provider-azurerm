package common

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func expandAzureRmCosmosDBIndexingPolicyIncludedPaths(input []interface{}) *[]documentdb.IncludedPath {
	if len(input) == 0 {
		return nil
	}

	var includedPaths []documentdb.IncludedPath

	for _, v := range input {
		includedPath := v.(map[string]interface{})
		path := documentdb.IncludedPath{
			Path: utils.String(includedPath["path"].(string)),
		}

		includedPaths = append(includedPaths, path)
	}

	return &includedPaths
}

func expandAzureRmCosmosDBIndexingPolicyExcludedPaths(input []interface{}) *[]documentdb.ExcludedPath {
	if len(input) == 0 {
		return nil
	}

	var paths []documentdb.ExcludedPath

	for _, v := range input {
		block := v.(map[string]interface{})
		paths = append(paths, documentdb.ExcludedPath{
			Path: utils.String(block["path"].(string)),
		})
	}

	return &paths
}

func ExpandAzureRmCosmosDBIndexingPolicyCompositeIndexes(input []interface{}) *[][]documentdb.CompositePath {
	indexes := make([][]documentdb.CompositePath, 0)

	for _, i := range input {
		indexPairs := make([]documentdb.CompositePath, 0)
		indexPair := i.(map[string]interface{})
		for _, idxPair := range indexPair["index"].([]interface{}) {
			data := idxPair.(map[string]interface{})

			index := documentdb.CompositePath{
				Path:  utils.String(data["path"].(string)),
				Order: documentdb.CompositePathSortOrder(strings.ToLower(data["order"].(string))),
			}
			indexPairs = append(indexPairs, index)
		}
		indexes = append(indexes, indexPairs)
	}

	return &indexes
}

func ExpandAzureRmCosmosDBIndexingPolicySpatialIndexes(input []interface{}) *[]documentdb.SpatialSpec {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	indexes := make([]documentdb.SpatialSpec, 0)
	// no matter what spatial types are updated, all types will be set and returned from service
	spatialTypes := []documentdb.SpatialType{
		documentdb.SpatialTypeLineString,
		documentdb.SpatialTypeMultiPolygon,
		documentdb.SpatialTypePoint,
		documentdb.SpatialTypePolygon,
	}

	for _, i := range input {
		indexPair := i.(map[string]interface{})
		indexes = append(indexes, documentdb.SpatialSpec{
			Types: &spatialTypes,
			Path:  utils.String(indexPair["path"].(string)),
		})
	}

	return &indexes
}

func ExpandAzureRmCosmosDbIndexingPolicy(d *pluginsdk.ResourceData) *documentdb.IndexingPolicy {
	i := d.Get("indexing_policy").([]interface{})

	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]interface{})
	policy := &documentdb.IndexingPolicy{}
	policy.IndexingMode = documentdb.IndexingMode(strings.ToLower(input["indexing_mode"].(string)))
	if v, ok := input["included_path"].([]interface{}); ok {
		policy.IncludedPaths = expandAzureRmCosmosDBIndexingPolicyIncludedPaths(v)
	}
	if v, ok := input["excluded_path"].([]interface{}); ok {
		policy.ExcludedPaths = expandAzureRmCosmosDBIndexingPolicyExcludedPaths(v)
	}

	if v, ok := input["composite_index"].([]interface{}); ok {
		policy.CompositeIndexes = ExpandAzureRmCosmosDBIndexingPolicyCompositeIndexes(v)
	}

	policy.SpatialIndexes = ExpandAzureRmCosmosDBIndexingPolicySpatialIndexes(input["spatial_index"].([]interface{}))

	return policy
}

func flattenCosmosDBIndexingPolicyExcludedPaths(input *[]documentdb.ExcludedPath) []interface{} {
	if input == nil {
		return nil
	}

	excludedPaths := make([]interface{}, 0)

	for _, v := range *input {
		// _etag is automatically added by the server and should be excluded on flattening
		// as the user isn't setting it and it will show changes in state.
		if *v.Path == "/\"_etag\"/?" {
			continue
		}

		block := make(map[string]interface{})
		block["path"] = v.Path
		excludedPaths = append(excludedPaths, block)
	}

	return excludedPaths
}

func flattenCosmosDBIndexingPolicyCompositeIndex(input []documentdb.CompositePath) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	indexPairs := make([]interface{}, 0)
	for _, v := range input {
		path := ""
		if v.Path != nil {
			path = *v.Path
		}

		block := make(map[string]interface{})
		block["path"] = path
		block["order"] = strings.Title(string(v.Order))
		indexPairs = append(indexPairs, block)
	}

	return indexPairs
}

func FlattenCosmosDBIndexingPolicyCompositeIndexes(input *[][]documentdb.CompositePath) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	indexes := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string][]interface{})
		block["index"] = flattenCosmosDBIndexingPolicyCompositeIndex(v)
		indexes = append(indexes, block)
	}

	return indexes
}

func flattenCosmosDBIndexingPolicyIncludedPaths(input *[]documentdb.IncludedPath) []interface{} {
	if input == nil {
		return nil
	}

	includedPaths := make([]interface{}, 0)

	for _, v := range *input {
		block := make(map[string]interface{})
		block["path"] = v.Path
		includedPaths = append(includedPaths, block)
	}

	return includedPaths
}

func FlattenCosmosDBIndexingPolicySpatialIndexes(input *[]documentdb.SpatialSpec) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	indexes := make([]interface{}, 0)

	for _, v := range *input {
		var path string
		if v.Path != nil {
			path = *v.Path
		}
		indexes = append(indexes, map[string]interface{}{
			"path":  path,
			"types": flattenCosmosDBIndexingPolicySpatialIndexesTypes(v.Types),
		})
	}

	return indexes
}

func flattenCosmosDBIndexingPolicySpatialIndexesTypes(input *[]documentdb.SpatialType) []interface{} {
	if input == nil {
		return nil
	}

	types := make([]interface{}, 0)

	for _, v := range *input {
		types = append(types, string(v))
	}

	return types
}

func FlattenAzureRmCosmosDbIndexingPolicy(indexingPolicy *documentdb.IndexingPolicy) []interface{} {
	results := make([]interface{}, 0)
	if indexingPolicy == nil {
		return results
	}

	result := make(map[string]interface{})
	result["indexing_mode"] = strings.Title(string(indexingPolicy.IndexingMode))
	result["included_path"] = flattenCosmosDBIndexingPolicyIncludedPaths(indexingPolicy.IncludedPaths)
	result["excluded_path"] = flattenCosmosDBIndexingPolicyExcludedPaths(indexingPolicy.ExcludedPaths)
	result["composite_index"] = FlattenCosmosDBIndexingPolicyCompositeIndexes(indexingPolicy.CompositeIndexes)
	result["spatial_index"] = FlattenCosmosDBIndexingPolicySpatialIndexes(indexingPolicy.SpatialIndexes)

	results = append(results, result)
	return results
}

func ValidateAzureRmCosmosDbIndexingPolicy(indexingPolicy *documentdb.IndexingPolicy) error {
	if indexingPolicy == nil {
		return nil
	}

	// Ensure includedPaths or excludedPaths are not set if indexingMode is "None".
	if indexingPolicy.IndexingMode == documentdb.None {
		if indexingPolicy.IncludedPaths != nil {
			return fmt.Errorf("included_path must not be set if indexing_mode is %q", strings.Title(string(documentdb.None)))
		}

		if indexingPolicy.ExcludedPaths != nil {
			return fmt.Errorf("excluded_path must not be set if indexing_mode is %q", strings.Title(string(documentdb.None)))
		}
	}

	// Any indexing policy has to include the root path /* as either an included or an excluded path.
	rootPath := "/*"
	includedPathsDefined := indexingPolicy.IncludedPaths != nil
	includedPathsContainRootPath := false

	if includedPathsDefined {
		for _, includedPath := range *indexingPolicy.IncludedPaths {
			if includedPathsContainRootPath {
				break
			}

			includedPathsContainRootPath = *includedPath.Path == rootPath
		}
	}

	excludedPathsContainRootPath := false
	excludedPathsDefined := indexingPolicy.ExcludedPaths != nil

	if excludedPathsDefined {
		for _, excludedPath := range *indexingPolicy.ExcludedPaths {
			if excludedPathsContainRootPath {
				break
			}

			excludedPathsContainRootPath = *excludedPath.Path == rootPath
		}
	}

	// The root path can't be included and excluded at the same time.
	if includedPathsContainRootPath && excludedPathsContainRootPath {
		return fmt.Errorf("only one of included_path or excluded_path may include the path %q", rootPath)
	}

	// The root path must be included or excluded
	if (includedPathsDefined || excludedPathsDefined) && !(includedPathsContainRootPath || excludedPathsContainRootPath) {
		return fmt.Errorf("either included_path or excluded_path must include the path %q", rootPath)
	}

	return nil
}
