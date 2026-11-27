// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandFusionSourceSettings(input []any) *[]alertrules.FusionSourceSettings {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSourceSettings, 0)

	for _, e := range input {
		e := e.(map[string]any)
		setting := alertrules.FusionSourceSettings{
			Enabled:        e["enabled"].(bool),
			SourceName:     e["name"].(string),
			SourceSubTypes: expandFusionSourceSubTypes(e["sub_type"].([]any)),
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSourceSubTypes(input []any) *[]alertrules.FusionSourceSubTypeSetting {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSourceSubTypeSetting, 0)

	for _, e := range input {
		e := e.(map[string]any)
		setting := alertrules.FusionSourceSubTypeSetting{
			Enabled:           e["enabled"].(bool),
			SourceSubTypeName: e["name"].(string),
			SeverityFilters: alertrules.FusionSubTypeSeverityFilter{
				Filters: expandFusionSubTypeSeverityFiltersItems(e["severities_allowed"].(*pluginsdk.Set).List()),
			},
		}
		result = append(result, setting)
	}

	return &result
}

func expandFusionSubTypeSeverityFiltersItems(input []any) *[]alertrules.FusionSubTypeSeverityFiltersItem {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FusionSubTypeSeverityFiltersItem, 0)

	// We can't simply remove the disabled properties in the request, as that will be reflected to the backend model (i.e. those unspecified severity will be absent also).
	// As any absent severity then will not be shown in the Portal when users try to edit the alert rule. The drop down menu won't show these absent severities...
	filters := map[string]bool{}
	for _, e := range alertrules.PossibleValuesForAlertSeverity() {
		filters[e] = false
	}

	for _, e := range input {
		filters[e.(string)] = true
	}

	for severity, enabled := range filters {
		item := alertrules.FusionSubTypeSeverityFiltersItem{
			Enabled:  enabled,
			Severity: alertrules.AlertSeverity(severity),
		}
		result = append(result, item)
	}

	return &result
}
