// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandMonitorMetricAlertCriteria(d *pluginsdk.ResourceData, isLegacy bool) (metricalerts.MetricAlertCriteria, error) {
	switch {
	case len(d.Get("criteria").([]any)) != 0:
		if isLegacy {
			return expandMonitorMetricAlertSingleResourceMultiMetricCriteria(d.Get("criteria").([]any)), nil
		}
		return expandMonitorMetricAlertMultiResourceMultiMetricForStaticMetricCriteria(d.Get("criteria").([]any)), nil
	case len(d.Get("dynamic_criteria").([]any)) != 0:
		return expandMonitorMetricAlertMultiResourceMultiMetricForDynamicMetricCriteria(d.Get("dynamic_criteria").([]any)), nil
	case len(d.Get("application_insights_web_test_location_availability_criteria").([]any)) != 0:
		return expandMonitorMetricAlertWebtestLocAvailCriteria(d.Get("application_insights_web_test_location_availability_criteria").([]any)), nil
	default:
		// Guaranteed by schema `AtLeastOne` constraint
		return nil, fmt.Errorf("unknown criteria type")
	}
}

func expandMonitorMetricAlertSingleResourceMultiMetricCriteria(input []any) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]any)
		dimensions := expandMonitorMetricDimension(v["dimension"].([]any))
		criteria = append(criteria, metricalerts.MetricCriteria{
			Name:                 fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:      pointer.To(v["metric_namespace"].(string)),
			MetricName:           v["metric_name"].(string),
			TimeAggregation:      metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:           &dimensions,
			Operator:             metricalerts.Operator(v["operator"].(string)),
			Threshold:            v["threshold"].(float64),
			SkipMetricValidation: pointer.To(v["skip_metric_validation"].(bool)),
		})
	}
	return &metricalerts.MetricAlertSingleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertMultiResourceMultiMetricForStaticMetricCriteria(input []any) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MultiMetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]any)
		dimensions := expandMonitorMetricDimension(v["dimension"].([]any))
		criteria = append(criteria, metricalerts.MetricCriteria{
			Name:                 fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:      pointer.To(v["metric_namespace"].(string)),
			MetricName:           v["metric_name"].(string),
			TimeAggregation:      metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:           &dimensions,
			Operator:             metricalerts.Operator(v["operator"].(string)),
			Threshold:            v["threshold"].(float64),
			SkipMetricValidation: pointer.To(v["skip_metric_validation"].(bool)),
		})
	}
	return &metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertMultiResourceMultiMetricForDynamicMetricCriteria(input []any) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MultiMetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]any)
		dimensions := expandMonitorMetricDimension(v["dimension"].([]any))

		dynamicMetricCriteria := metricalerts.DynamicMetricCriteria{
			Name:             fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:  pointer.To(v["metric_namespace"].(string)),
			MetricName:       v["metric_name"].(string),
			TimeAggregation:  metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:       &dimensions,
			Operator:         metricalerts.DynamicThresholdOperator(v["operator"].(string)),
			AlertSensitivity: metricalerts.DynamicThresholdSensitivity(v["alert_sensitivity"].(string)),
			FailingPeriods: metricalerts.DynamicThresholdFailingPeriods{
				NumberOfEvaluationPeriods: float64(v["evaluation_total_count"].(int)),
				MinFailingPeriodsToAlert:  float64(v["evaluation_failure_count"].(int)),
			},
			SkipMetricValidation: pointer.To(v["skip_metric_validation"].(bool)),
		}

		if datetime := v["ignore_data_before"].(string); datetime != "" {
			// Guaranteed in schema validation func.
			t, _ := time.Parse(time.RFC3339, datetime)
			ignoreDataBefore := &date.Time{Time: t}
			dynamicMetricCriteria.IgnoreDataBefore = pointer.To(ignoreDataBefore.String())
		}

		criteria = append(criteria, dynamicMetricCriteria)
	}
	return &metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertWebtestLocAvailCriteria(input []any) metricalerts.MetricAlertCriteria {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &metricalerts.WebtestLocationAvailabilityCriteria{
		WebTestId:           v["web_test_id"].(string),
		ComponentId:         v["component_id"].(string),
		FailedLocationCount: float64(v["failed_location_count"].(int)),
	}
}

func expandMonitorMetricDimension(input []any) []metricalerts.MetricDimension {
	result := make([]metricalerts.MetricDimension, 0)
	for _, dimension := range input {
		dVal := dimension.(map[string]any)
		result = append(result, metricalerts.MetricDimension{
			Name:     dVal["name"].(string),
			Operator: dVal["operator"].(string),
			Values:   helpers.ExpandStringValues(dVal["values"].([]any)),
		})
	}
	return result
}

func expandMonitorMetricAlertAction(input []any) *[]metricalerts.MetricAlertAction {
	actions := make([]metricalerts.MetricAlertAction, 0)
	for _, item := range input {
		v := item.(map[string]any)
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]any) {
					props[pk] = pv.(string)
				}
			}

			actions = append(actions, metricalerts.MetricAlertAction{
				ActionGroupId:     pointer.To(agID),
				WebHookProperties: &props,
			})
		}
	}
	return &actions
}
