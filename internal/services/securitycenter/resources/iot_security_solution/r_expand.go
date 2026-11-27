// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_solution

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandIotSecuritySolutionExport(input []any) *[]security.ExportData {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	result := make([]security.ExportData, 0)
	for _, item := range input {
		result = append(result, security.ExportData(item.(string)))
	}
	return &result
}

func expandIotSecuritySolutionRecommendation(input []any) *[]security.RecommendationConfigurationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	result := make([]security.RecommendationConfigurationProperties, 0)
	v := input[0].(map[string]any)
	for k, item := range getRecommendationSchemaMap() {
		status := security.Disabled
		if v[item].(bool) {
			status = security.Enabled
		}
		result = append(result, security.RecommendationConfigurationProperties{
			RecommendationType: k,
			Status:             status,
		})
	}
	return &result
}

func expandIotSecuritySolutionAdditionalWorkspace(input []any) *[]security.AdditionalWorkspacesProperties {
	results := make([]security.AdditionalWorkspacesProperties, 0)

	for _, item := range input {
		v := item.(map[string]any)

		dataTypes := make([]security.AdditionalWorkspaceDataType, 0)
		for _, item := range *pluginsdk.ExpandStringSlice(v["data_types"].(*pluginsdk.Set).List()) {
			dataTypes = append(dataTypes, security.AdditionalWorkspaceDataType(item))
		}

		results = append(results, security.AdditionalWorkspacesProperties{
			Workspace: pointer.To(v["workspace_id"].(string)),
			Type:      security.Sentinel,
			DataTypes: &dataTypes,
		})
	}

	return &results
}

func expandIotSecuritySolutionDisabledDataSources(input []any) *[]security.DataSource {
	if len(input) == 0 {
		return nil
	}

	disabledDataSources := make([]security.DataSource, 0)
	for _, item := range *pluginsdk.ExpandStringSlice(input) {
		disabledDataSources = append(disabledDataSources, security.DataSource(item))
	}

	return &disabledDataSources
}
