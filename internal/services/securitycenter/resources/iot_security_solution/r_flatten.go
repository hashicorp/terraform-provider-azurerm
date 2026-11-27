// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_solution

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenIotSecuritySolutionExport(input *[]security.ExportData) []any {
	result := make([]any, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenIotSecuritySolutionRecommendation(input *[]security.RecommendationConfigurationProperties) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}
	result := make(map[string]any)
	schemaMap := getRecommendationSchemaMap()
	for _, item := range *input {
		if v, ok := schemaMap[item.RecommendationType]; ok {
			result[v] = item.Status == security.Enabled
		}
	}
	return []any{result}
}

func flattenIotSecuritySolutionAdditionalWorkspace(input *[]security.AdditionalWorkspacesProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		rawDataTypes := make([]string, 0)
		for _, item := range *item.DataTypes {
			rawDataTypes = append(rawDataTypes, string(item))
		}
		dataTypes := pluginsdk.FlattenSlice(&rawDataTypes)

		results = append(results, map[string]any{
			"data_types":   dataTypes,
			"workspace_id": pointer.From(item.Workspace),
		})
	}

	return results
}

func flattenIotSecuritySolutionDisabledDataSources(input *[]security.DataSource) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return pluginsdk.FlattenSlice(&results)
}
