// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"fmt"
	"maps"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMonitorMetricAlertFlatten(d *pluginsdk.ResourceData, id *metricalerts.MetricAlertId, model *metricalerts.MetricAlertResource) error {
	d.Set("name", id.MetricAlertName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		props := model.Properties
		d.Set("enabled", props.Enabled)
		d.Set("auto_mitigate", props.AutoMitigate)
		d.Set("description", props.Description)
		d.Set("severity", props.Severity)
		d.Set("frequency", props.EvaluationFrequency)
		d.Set("window_size", props.WindowSize)
		if err := d.Set("scopes", props.Scopes); err != nil {
			return fmt.Errorf("setting `scopes`: %+v", err)
		}

		// Determine the correct criteria schema to set
		var criteriaSchema string
		switch c := props.Criteria.(type) {
		case metricalerts.MetricAlertSingleResourceMultipleMetricCriteria:
			criteriaSchema = "criteria"
		case metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria:
			if c.AllOf == nil || len(*c.AllOf) == 0 {
				return fmt.Errorf("nil or empty contained criteria of MultipleResourceMultipleMetricCriteria")
			}
			// `MinItems` defined in schema guaranteed there is at least one element.
			switch (*c.AllOf)[0].(type) {
			case metricalerts.DynamicMetricCriteria:
				criteriaSchema = "dynamic_criteria"
			case metricalerts.MetricCriteria:
				criteriaSchema = "criteria"
			}
		case metricalerts.WebtestLocationAvailabilityCriteria:
			criteriaSchema = "application_insights_web_test_location_availability_criteria"
		default:
			return fmt.Errorf("unknown criteria type")
		}

		monitorMetricAlertCriteria := flattenMonitorMetricAlertCriteria(props.Criteria)
		// lintignore:R001
		if err := d.Set(criteriaSchema, monitorMetricAlertCriteria); err != nil {
			return fmt.Errorf("failed setting `%s`: %+v", criteriaSchema, err)
		}

		if err := d.Set("action", flattenMonitorMetricAlertAction(props.Actions)); err != nil {
			return fmt.Errorf("setting `action`: %+v", err)
		}
		d.Set("target_resource_type", props.TargetResourceType)
		d.Set("target_resource_location", props.TargetResourceRegion)

		if err := d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMonitorMetricAlertCriteria(input metricalerts.MetricAlertCriteria) []any {
	switch criteria := input.(type) {
	case metricalerts.MetricAlertSingleResourceMultipleMetricCriteria:
		return flattenMonitorMetricAlertSingleResourceMultiMetricCriteria(criteria.AllOf)
	case metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria:
		return flattenMonitorMetricAlertMultiResourceMultiMetricCriteria(criteria.AllOf)
	case metricalerts.WebtestLocationAvailabilityCriteria:
		return flattenMonitorMetricAlertWebtestLocAvailCriteria(&criteria)
	default:
		return []any{}
	}
}

func flattenMonitorMetricAlertSingleResourceMultiMetricCriteria(input *[]metricalerts.MetricCriteria) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}
	criteria := (*input)[0]
	metricName := criteria.MetricName
	metricNamespace := criteria.MetricNamespace
	timeAggregation := criteria.TimeAggregation

	dimResult := make([]map[string]any, 0)
	if criteria.Dimensions != nil {
		for _, dimension := range *criteria.Dimensions {
			dVal := make(map[string]any)
			dVal["name"] = dimension.Name
			dVal["operator"] = dimension.Operator
			dVal["values"] = dimension.Values
			dimResult = append(dimResult, dVal)
		}
	}

	operator := string(criteria.Operator)
	threshold := criteria.Threshold

	skipMetricValidation := pointer.From(criteria.SkipMetricValidation)

	return []any{
		map[string]any{
			"metric_namespace":       metricNamespace,
			"metric_name":            metricName,
			"aggregation":            timeAggregation,
			"dimension":              dimResult,
			"operator":               operator,
			"threshold":              threshold,
			"skip_metric_validation": skipMetricValidation,
		},
	}
}

func flattenMonitorMetricAlertMultiResourceMultiMetricCriteria(input *[]metricalerts.MultiMetricCriteria) []any {
	if input == nil {
		return []any{}
	}
	result := make([]any, 0)

	for _, criteria := range *input {
		v := make(map[string]any)
		var (
			metricName           string
			metricNamespace      string
			timeAggregation      any
			dimensions           []metricalerts.MetricDimension
			skipMetricValidation bool
		)

		switch criteria := criteria.(type) {
		case metricalerts.MetricCriteria:
			metricName = criteria.MetricName

			if criteria.MetricNamespace != nil {
				metricNamespace = *criteria.MetricNamespace
			}
			timeAggregation = criteria.TimeAggregation
			if criteria.Dimensions != nil {
				dimensions = *criteria.Dimensions
			}

			// MetricCriteria specific properties
			v["operator"] = string(criteria.Operator)
			v["threshold"] = criteria.Threshold
			if criteria.SkipMetricValidation != nil {
				skipMetricValidation = *criteria.SkipMetricValidation
			}
		case metricalerts.DynamicMetricCriteria:
			metricName = criteria.MetricName

			if criteria.MetricNamespace != nil {
				metricNamespace = *criteria.MetricNamespace
			}
			timeAggregation = criteria.TimeAggregation
			if criteria.Dimensions != nil {
				dimensions = *criteria.Dimensions
			}
			if criteria.SkipMetricValidation != nil {
				skipMetricValidation = *criteria.SkipMetricValidation
			}
			// DynamicMetricCriteria specific properties
			v["operator"] = string(criteria.Operator)
			v["alert_sensitivity"] = string(criteria.AlertSensitivity)

			v["evaluation_total_count"] = int(criteria.FailingPeriods.NumberOfEvaluationPeriods)
			v["evaluation_failure_count"] = int(criteria.FailingPeriods.MinFailingPeriodsToAlert)

			v["ignore_data_before"] = pointer.From(criteria.IgnoreDataBefore)
		}

		// Common properties
		v["metric_name"] = metricName
		v["metric_namespace"] = metricNamespace
		v["aggregation"] = timeAggregation
		v["skip_metric_validation"] = skipMetricValidation
		if dimensions != nil {
			dimResult := make([]map[string]any, 0)
			for _, dimension := range dimensions {
				dVal := make(map[string]any)
				dVal["name"] = dimension.Name
				dVal["operator"] = dimension.Operator
				dVal["values"] = dimension.Values
				dimResult = append(dimResult, dVal)
			}
			v["dimension"] = dimResult
		}

		result = append(result, v)
	}
	return result
}

func flattenMonitorMetricAlertWebtestLocAvailCriteria(input *metricalerts.WebtestLocationAvailabilityCriteria) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"web_test_id":           input.WebTestId,
			"component_id":          input.ComponentId,
			"failed_location_count": int(input.FailedLocationCount),
		},
	}
}

func flattenMonitorMetricAlertAction(input *[]metricalerts.MetricAlertAction) (result []any) {
	result = make([]any, 0)
	if input == nil {
		return
	}
	for _, action := range *input {
		v := make(map[string]any)

		if action.ActionGroupId != nil {
			v["action_group_id"] = *action.ActionGroupId
		}

		props := make(map[string]string)
		if action.WebHookProperties != nil {
			maps.Copy(props, *action.WebHookProperties)
		}
		v["webhook_properties"] = props

		result = append(result, v)
	}

	return result
}
