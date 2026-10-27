// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
)

func flattenLogRules(input *tagrules.LogRules) []any {
	results := make([]any, 0)

	if input != nil {
		results = append(results, map[string]any{
			"aad_log_enabled":          pointer.From(input.SendAadLogs),
			"filter":                   flattenFilteringTags(input.FilteringTags),
			"resource_log_enabled":     pointer.From(input.SendResourceLogs),
			"subscription_log_enabled": pointer.From(input.SendSubscriptionLogs),
		})
	}

	return results
}

func flattenMetricRules(input *tagrules.MetricRules) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"filter": flattenFilteringTags(input.FilteringTags),
		},
	}
}

func flattenFilteringTags(input *[]tagrules.FilteringTag) []any {
	results := make([]any, 0)
	if input != nil {
		for _, filteringTagRules := range *input {
			action := ""
			if filteringTagRules.Action != nil {
				action = string(*filteringTagRules.Action)
			}
			results = append(results, map[string]any{
				"action": action,
				"name":   pointer.From(filteringTagRules.Name),
				"value":  pointer.From(filteringTagRules.Value),
			})
		}
	}
	return results
}
