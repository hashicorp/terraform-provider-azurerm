// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
)

func flattenFusionSourceSettings(input *[]alertrules.FusionSourceSettings) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, e := range *input {
		output = append(output, map[string]any{
			"name":     e.SourceName,
			"enabled":  e.Enabled,
			"sub_type": flattenFusionSourceSubTypes(e.SourceSubTypes),
		})
	}

	return output
}

func flattenFusionSourceSubTypes(input *[]alertrules.FusionSourceSubTypeSetting) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, e := range *input {
		output = append(output, map[string]any{
			"name":               e.SourceSubTypeName,
			"enabled":            e.Enabled,
			"severities_allowed": flattenFusionSubTypeSeverityFiltersItems(e.SeverityFilters.Filters),
		})
	}

	return output
}

func flattenFusionSubTypeSeverityFiltersItems(input *[]alertrules.FusionSubTypeSeverityFiltersItem) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, e := range *input {
		if e.Enabled {
			output = append(output, string(e.Severity))
		}
	}

	return output
}
