// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
)

func expandLogRules(input []any) *tagrules.LogRules {
	if len(input) == 0 {
		return nil
	}
	// an empty `log {}` block is returned as a nil element since every attribute is optional
	v, _ := input[0].(map[string]any)
	aadLogEnabled, _ := v["aad_log_enabled"].(bool)
	subscriptionLogEnabled, _ := v["subscription_log_enabled"].(bool)
	resourceLogEnabled, _ := v["resource_log_enabled"].(bool)
	filteringTag, _ := v["filter"].([]any)

	return &tagrules.LogRules{
		SendAadLogs:          pointer.To(aadLogEnabled),
		SendSubscriptionLogs: pointer.To(subscriptionLogEnabled),
		SendResourceLogs:     pointer.To(resourceLogEnabled),
		FilteringTags:        expandFilteringTag(filteringTag),
	}
}

func expandMetricRules(input []any) *tagrules.MetricRules {
	if len(input) == 0 {
		return nil
	}
	// an empty `metric {}` block is returned as a nil element since every attribute is optional
	v, _ := input[0].(map[string]any)
	filteringTag, _ := v["filter"].([]any)

	return &tagrules.MetricRules{
		FilteringTags: expandFilteringTag(filteringTag),
	}
}

func expandFilteringTag(input []any) *[]tagrules.FilteringTag {
	filteringTags := make([]tagrules.FilteringTag, 0)

	for _, v := range input {
		config := v.(map[string]any)

		filteringTags = append(filteringTags, tagrules.FilteringTag{
			Name:   pointer.To(config["name"].(string)),
			Value:  pointer.To(config["value"].(string)),
			Action: pointer.ToEnum[tagrules.TagAction](config["action"].(string)),
		})
	}

	return &filteringTags
}
