package monitor

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func expandMonitorScheduledQueryRulesCommonSource(d *schema.ResourceData) *insights.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*schema.Set).List()
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

// ValidateThreshold checks that a threshold value is between 0 and 10000
// and is a whole number. The azure-sdk-for-go expects this value to be a float64
// but the user validation rules want an integer.
func validateMonitorScheduledQueryRulesAlertThreshold(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float64", k))
	}

	if v != float64(int64(v)) {
		errors = append(errors, fmt.Errorf("%q must be a whole number", k))
	}

	if v < 0 || v > 10000 {
		errors = append(errors, fmt.Errorf("%q must be between 0 and 10000 (inclusive)", k))
	}

	return warnings, errors
}
