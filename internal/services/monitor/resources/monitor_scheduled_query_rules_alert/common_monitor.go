// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
)

func flattenAzureRmScheduledQueryRulesAlertAction(input *scheduledqueryrules.AzNsActionGroup) []any {
	v := make(map[string]any)

	if input != nil {
		if input.ActionGroup != nil {
			v["action_group"] = *input.ActionGroup
		}
		v["email_subject"] = input.EmailSubject
		v["custom_webhook_payload"] = input.CustomWebhookPayload
	}
	return []any{v}
}

func flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input *scheduledqueryrules.LogMetricTrigger) []any {
	result := make(map[string]any)

	if input == nil {
		return []any{}
	}

	result["operator"] = input.ThresholdOperator

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	result["metric_trigger_type"] = input.MetricTriggerType

	if input.MetricColumn != nil {
		result["metric_column"] = *input.MetricColumn
	}
	return []any{result}
}

func flattenAzureRmScheduledQueryRulesAlertTrigger(input scheduledqueryrules.TriggerCondition) []any {
	result := make(map[string]any)

	result["operator"] = string(input.ThresholdOperator)
	result["threshold"] = input.Threshold

	if input.MetricTrigger != nil {
		result["metric_trigger"] = flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input.MetricTrigger)
	}

	return []any{result}
}
