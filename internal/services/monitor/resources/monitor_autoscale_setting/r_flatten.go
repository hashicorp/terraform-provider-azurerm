// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_autoscale_setting

import (
	"fmt"
	"maps"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
)

func flattenAzureRmMonitorAutoScaleSettingProfile(profiles []autoscalesettings.AutoscaleProfile) ([]any, error) {
	if profiles == nil {
		return []any{}, nil
	}

	results := make([]any, 0)
	for _, profile := range profiles {
		result := make(map[string]any)

		result["name"] = profile.Name

		capacity, err := flattenAzureRmMonitorAutoScaleSettingCapacity(profile.Capacity)
		if err != nil {
			return nil, fmt.Errorf("flattening `capacity`: %+v", err)
		}
		result["capacity"] = capacity

		result["fixed_date"] = flattenAzureRmMonitorAutoScaleSettingFixedDate(profile.FixedDate)
		result["recurrence"] = flattenAzureRmMonitorAutoScaleSettingRecurrence(profile.Recurrence)

		rule, err := flattenAzureRmMonitorAutoScaleSettingRules(profile.Rules)
		if err != nil {
			return nil, fmt.Errorf("flattening Rule: %s", err)
		}
		result["rule"] = rule

		results = append(results, result)
	}
	return results, nil
}

func flattenAzureRmMonitorAutoScaleSettingPredictive(input *autoscalesettings.PredictiveAutoscalePolicy) []any {
	if input == nil {
		return []any{}
	}

	// omit the block if disabled
	if input.ScaleMode == autoscalesettings.PredictiveAutoscalePolicyScaleModeDisabled {
		return []any{}
	}

	result := map[string]any{
		"look_ahead_time": pointer.From(input.ScaleLookAheadTime),
		"scale_mode":      string(input.ScaleMode),
	}

	return []any{result}
}

func flattenAzureRmMonitorAutoScaleSettingCapacity(input autoscalesettings.ScaleCapacity) ([]any, error) {
	result := make(map[string]any)

	min, err := strconv.Atoi(input.Minimum)
	if err != nil {
		return nil, fmt.Errorf("converting Minimum Scale Capacity %q to an int: %+v", input.Minimum, err)
	}
	result["minimum"] = min

	max, err := strconv.Atoi(input.Maximum)
	if err != nil {
		return nil, fmt.Errorf("converting Maximum Scale Capacity %q to an int: %+v", input.Maximum, err)
	}
	result["maximum"] = max

	defaultCapacity, err := strconv.Atoi(input.Default)
	if err != nil {
		return nil, fmt.Errorf("converting Default Scale Capacity %q to an int: %+v", input.Default, err)
	}
	result["default"] = defaultCapacity

	return []any{result}, nil
}

func flattenAzureRmMonitorAutoScaleSettingRules(input []autoscalesettings.ScaleRule) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	results := make([]any, 0)
	for _, rule := range input {
		result := make(map[string]any)

		metricTriggers := make([]any, 0)
		var metricNamespace string
		var dividePerInstance bool

		if v := rule.MetricTrigger.MetricNamespace; v != nil {
			metricNamespace = *v
		}

		if rule.MetricTrigger.DividePerInstance != nil {
			dividePerInstance = *rule.MetricTrigger.DividePerInstance
		}

		metricTriggers = append(metricTriggers, map[string]any{
			"metric_name":              rule.MetricTrigger.MetricName,
			"metric_namespace":         metricNamespace,
			"metric_resource_id":       rule.MetricTrigger.MetricResourceUri,
			"time_grain":               rule.MetricTrigger.TimeGrain,
			"statistic":                string(rule.MetricTrigger.Statistic),
			"time_window":              rule.MetricTrigger.TimeWindow,
			"time_aggregation":         string(rule.MetricTrigger.TimeAggregation),
			"operator":                 string(rule.MetricTrigger.Operator),
			"threshold":                rule.MetricTrigger.Threshold,
			"dimensions":               flattenAzureRmMonitorAutoScaleSettingRulesDimensions(rule.MetricTrigger.Dimensions),
			"divide_by_instance_count": dividePerInstance,
		})

		result["metric_trigger"] = metricTriggers

		scaleActions := make([]any, 0)
		v := rule.ScaleAction
		action := make(map[string]any)

		action["direction"] = string(v.Direction)
		action["type"] = string(v.Type)
		action["cooldown"] = v.Cooldown

		if val := v.Value; val != nil && *val != "" {
			i, err := strconv.Atoi(*val)
			if err != nil {
				return nil, fmt.Errorf("`value` %q was not convertible to an int: %s", *val, err)
			}
			action["value"] = i
		}

		scaleActions = append(scaleActions, action)

		result["scale_action"] = scaleActions

		results = append(results, result)
	}

	return results, nil
}

func flattenAzureRmMonitorAutoScaleSettingFixedDate(input *autoscalesettings.TimeWindow) []any {
	if input == nil {
		return []any{}
	}

	result := make(map[string]any)

	if timezone := input.TimeZone; timezone != nil {
		result["timezone"] = *timezone
	}
	result["start"] = input.Start
	result["end"] = input.End

	return []any{result}
}

func flattenAzureRmMonitorAutoScaleSettingRecurrence(input *autoscalesettings.Recurrence) []any {
	if input == nil {
		return []any{}
	}

	result := make(map[string]any)

	schedule := input.Schedule
	result["timezone"] = schedule.TimeZone

	days := make([]string, 0)
	if s := schedule.Days; s != nil {
		days = s
	}
	result["days"] = days

	hours := make([]int, 0)
	if schedule.Hours != nil {
		for _, v := range schedule.Hours {
			hours = append(hours, int(v))
		}
	}
	result["hours"] = hours

	minutes := make([]int, 0)
	if schedule.Minutes != nil {
		for _, v := range schedule.Minutes {
			minutes = append(minutes, int(v))
		}
	}
	result["minutes"] = minutes

	return []any{result}
}

func flattenAzureRmMonitorAutoScaleSettingNotification(notifications *[]autoscalesettings.AutoscaleNotification) []any {
	results := make([]any, 0)

	if notifications == nil {
		return results
	}

	for _, notification := range *notifications {
		result := make(map[string]any)

		emails := make([]any, 0)
		if email := notification.Email; email != nil {
			block := make(map[string]any)

			if send := email.SendToSubscriptionAdministrator; send != nil {
				block["send_to_subscription_administrator"] = *send
			}

			if send := email.SendToSubscriptionCoAdministrators; send != nil {
				block["send_to_subscription_co_administrator"] = *send
			}

			customEmails := make([]any, 0)
			if custom := email.CustomEmails; custom != nil {
				for _, v := range *custom {
					customEmails = append(customEmails, v)
				}
			}
			block["custom_emails"] = customEmails

			emails = append(emails, block)
		}
		result["email"] = emails

		webhooks := make([]any, 0)
		if hooks := notification.WebHooks; hooks != nil {
			for _, v := range *hooks {
				hook := make(map[string]any)

				if v.ServiceUri != nil {
					hook["service_uri"] = *v.ServiceUri
				}

				props := make(map[string]string)
				if webHookProps := v.Properties; webHookProps != nil {
					maps.Copy(props, *v.Properties)
					hook["properties"] = props
					webhooks = append(webhooks, hook)
				}
			}
		}

		result["webhook"] = webhooks

		results = append(results, result)
	}
	return results
}

func flattenAzureRmMonitorAutoScaleSettingRulesDimensions(dimensions *[]autoscalesettings.ScaleRuleMetricDimension) []any {
	results := make([]any, 0)

	if dimensions == nil {
		return results
	}

	for _, dimension := range *dimensions {
		results = append(results, map[string]any{
			"name":     dimension.DimensionName,
			"operator": string(dimension.Operator),
			"values":   dimension.Values,
		})
	}
	return results
}
