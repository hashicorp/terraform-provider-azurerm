// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
)

func flattenOptimizedAutoScale(optimizedAutoScale *clusters.OptimizedAutoscale) []any {
	if optimizedAutoScale == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"maximum_instances": int(optimizedAutoScale.Maximum),
			"minimum_instances": int(optimizedAutoScale.Minimum),
		},
	}
}

func flattenKustoClusterSku(sku *clusters.AzureSku) []any {
	if sku == nil {
		return []any{}
	}

	s := map[string]any{
		"name": string(sku.Name),
	}

	if sku.Capacity != nil {
		s["capacity"] = int(*sku.Capacity)
	}

	return []any{s}
}

func flattenKustoClusterLanguageExtensionList(extensions *clusters.LanguageExtensionsList) []any {
	if extensions == nil {
		return []any{}
	}

	output := make([]any, 0)
	if extensions.Value != nil {
		for _, v := range *extensions.Value {
			output = append(
				output,
				map[string]any{
					"name":  string(*v.LanguageExtensionName),
					"image": string(*v.LanguageExtensionImageName),
				},
			)
		}
	}

	return output
}
