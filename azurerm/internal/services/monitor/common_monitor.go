package monitor

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func flattenAzureRmScheduledQueryRulesAlertAction(input *insights.AzNsActionGroup) []interface{} {
	v := make(map[string]interface{})

	if input != nil {
		if input.ActionGroup != nil {
			v["action_group"] = *input.ActionGroup
		}
		v["email_subject"] = input.EmailSubject
		v["custom_webhook_payload"] = input.CustomWebhookPayload
	}
	return []interface{}{v}
}

func expandMonitorScheduledQueryRulesCommonSource(d *pluginsdk.ResourceData) *insights.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*pluginsdk.Set).List()
	dataSourceID := d.Get("data_source_id").(string)
	query, ok := d.GetOk("query")
	source := insights.Source{
		AuthorizedResources: utils.ExpandStringSlice(authorizedResourceIDs),
		DataSourceID:        utils.String(dataSourceID),
		QueryType:           insights.ResultCount,
	}
	if ok {
		source.Query = utils.String(query.(string))
	}

	return &source
}

func flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input *insights.LogMetricTrigger) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	result["operator"] = string(input.ThresholdOperator)

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	result["metric_trigger_type"] = string(input.MetricTriggerType)

	if input.MetricColumn != nil {
		result["metric_column"] = *input.MetricColumn
	}
	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesAlertTrigger(input *insights.TriggerCondition) []interface{} {
	result := make(map[string]interface{})

	result["operator"] = string(input.ThresholdOperator)

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	if input.MetricTrigger != nil {
		result["metric_trigger"] = flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input.MetricTrigger)
	}

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesLogCriteria(input *[]insights.Criteria) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, criteria := range *input {
			v := make(map[string]interface{})

			v["dimension"] = flattenAzureRmScheduledQueryRulesLogDimension(criteria.Dimensions)
			v["metric_name"] = *criteria.MetricName

			result = append(result, v)
		}
	}

	return result
}

func flattenAzureRmScheduledQueryRulesLogDimension(input *[]insights.Dimension) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, dimension := range *input {
			v := make(map[string]interface{})

			if dimension.Name != nil {
				v["name"] = *dimension.Name
			}

			if dimension.Operator != nil {
				v["operator"] = *dimension.Operator
			}

			if dimension.Values != nil {
				v["values"] = *dimension.Values
			}

			result = append(result, v)
		}
	}

	return result
}
