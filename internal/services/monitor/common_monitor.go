// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

func expandMonitorScheduledQueryRulesCommonSource(d *pluginsdk.ResourceData) scheduledqueryrules.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*pluginsdk.Set).List()
	dataSourceID := d.Get("data_source_id").(string)

	source := scheduledqueryrules.Source{
		AuthorizedResources: pluginsdk.ExpandStringSlice(authorizedResourceIDs),
		DataSourceId:        dataSourceID,
	}

	if query, ok := d.GetOk("query"); ok {
		source.Query = pointer.To(query.(string))
	}
	if queryType, ok := d.GetOk("query_type"); ok {
		source.QueryType = pointer.ToEnum[scheduledqueryrules.QueryType](queryType.(string))
	}

	return source
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

func flattenAzureRmScheduledQueryRulesLogCriteria(input []scheduledqueryrules.Criteria) []any {
	result := make([]any, 0)
	for _, criteria := range input {
		v := make(map[string]any)

		v["dimension"] = flattenAzureRmScheduledQueryRulesLogDimension(criteria.Dimensions)
		v["metric_name"] = criteria.MetricName

		result = append(result, v)
	}
	return result
}

func flattenAzureRmScheduledQueryRulesLogDimension(input *[]scheduledqueryrules.Dimension) []any {
	result := make([]any, 0)

	if input != nil {
		for _, dimension := range *input {
			v := make(map[string]any)

			v["name"] = dimension.Name
			v["operator"] = dimension.Operator
			v["values"] = dimension.Values
			result = append(result, v)
		}
	}
	return result
}

func expandStringValues(input []any) []string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return result
}
