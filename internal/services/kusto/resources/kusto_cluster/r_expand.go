// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
)

func expandOptimizedAutoScale(input []any) *clusters.OptimizedAutoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]any)
	return &clusters.OptimizedAutoscale{
		Version:   1,
		IsEnabled: true,
		Minimum:   int64(config["minimum_instances"].(int)),
		Maximum:   int64(config["maximum_instances"].(int)),
	}
}

func expandKustoListString(input []any) *[]string {
	result := make([]string, 0)

	for _, v := range input {
		result = append(result, v.(string))
	}

	return &result
}

func expandKustoClusterSku(input []any) (*clusters.AzureSku, error) {
	sku := input[0].(map[string]any)
	name := sku["name"].(string)

	skuNamePrefixToTier := map[string]string{
		"Dev(No SLA)": "Basic",
		"Standard":    "Standard",
	}

	skuNamePrefix, _, _ := strings.Cut(sku["name"].(string), "_")
	tier, ok := skuNamePrefixToTier[skuNamePrefix]
	if !ok {
		return nil, fmt.Errorf("sku name begins with invalid tier, possible are Dev(No SLA) and Standard but is: %q", skuNamePrefix)
	}
	capacity := sku["capacity"].(int)

	azureSku := clusters.AzureSku{
		Name:     clusters.AzureSkuName(name),
		Tier:     clusters.AzureSkuTier(tier),
		Capacity: pointer.To(int64(capacity)),
	}

	return &azureSku, nil
}

func expandKustoClusterLanguageExtensionList(input []any) *clusters.LanguageExtensionsList {
	extensions := make([]clusters.LanguageExtension, 0)

	for _, ext := range input {
		extMap := ext.(map[string]any)
		extensions = append(extensions, clusters.LanguageExtension{
			LanguageExtensionName:      pointer.ToEnum[clusters.LanguageExtensionName](extMap["name"].(string)),
			LanguageExtensionImageName: pointer.ToEnum[clusters.LanguageExtensionImageName](extMap["image"].(string)),
		})
	}

	return &clusters.LanguageExtensionsList{
		Value: &extensions,
	}
}
