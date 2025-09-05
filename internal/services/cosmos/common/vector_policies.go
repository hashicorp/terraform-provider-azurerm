// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-04-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ExpandCosmosDbVectorEmbeddingPolicy(input []interface{}) *cosmosdb.VectorEmbeddingPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	policy := &cosmosdb.VectorEmbeddingPolicy{}

	if vectorEmbeddings, ok := v["vector_embedding"].([]interface{}); ok && len(vectorEmbeddings) > 0 {
		embeddings := make([]cosmosdb.VectorEmbedding, 0)
		for _, embedding := range vectorEmbeddings {
			e := embedding.(map[string]interface{})
			dataType := cosmosdb.VectorDataType(e["data_type"].(string))
			distanceFunction := cosmosdb.DistanceFunction(e["distance_function"].(string))
			embeddings = append(embeddings, cosmosdb.VectorEmbedding{
				Path:             e["path"].(string),
				DataType:         dataType,
				DistanceFunction: distanceFunction,
				Dimensions:       int64(e["dimensions"].(int)),
			})
		}
		policy.VectorEmbeddings = &embeddings
	}

	return policy
}

func FlattenCosmosDbVectorEmbeddingPolicy(policy *cosmosdb.VectorEmbeddingPolicy) []interface{} {
	if policy == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if policy.VectorEmbeddings != nil {
		embeddings := make([]interface{}, 0)
		for _, embedding := range *policy.VectorEmbeddings {
			e := make(map[string]interface{})
			e["path"] = embedding.Path
			e["data_type"] = string(embedding.DataType)
			e["distance_function"] = string(embedding.DistanceFunction)
			e["dimensions"] = int(embedding.Dimensions)
			embeddings = append(embeddings, e)
		}
		result["vector_embedding"] = embeddings
	}

	return []interface{}{result}
}

func ExpandCosmosDbFullTextPolicy(input []interface{}) *cosmosdb.FullTextPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	policy := &cosmosdb.FullTextPolicy{}

	if defaultLanguage, ok := v["default_language"].(string); ok && defaultLanguage != "" {
		policy.DefaultLanguage = utils.String(defaultLanguage)
	}

	if fullTextPaths, ok := v["full_text_path"].([]interface{}); ok && len(fullTextPaths) > 0 {
		paths := make([]cosmosdb.FullTextPath, 0)
		for _, pathInterface := range fullTextPaths {
			p := pathInterface.(map[string]interface{})
			path := cosmosdb.FullTextPath{
				Path: p["path"].(string),
			}
			if language, ok := p["language"].(string); ok && language != "" {
				path.Language = utils.String(language)
			}
			paths = append(paths, path)
		}
		policy.FullTextPaths = &paths
	}

	return policy
}

func FlattenCosmosDbFullTextPolicy(policy *cosmosdb.FullTextPolicy) []interface{} {
	if policy == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if policy.DefaultLanguage != nil {
		result["default_language"] = *policy.DefaultLanguage
	}

	if policy.FullTextPaths != nil {
		paths := make([]interface{}, 0)
		for _, path := range *policy.FullTextPaths {
			p := make(map[string]interface{})
			p["path"] = path.Path
			if path.Language != nil {
				p["language"] = *path.Language
			}
			paths = append(paths, p)
		}
		result["full_text_path"] = paths
	}

	return []interface{}{result}
}
