// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_autoscale_setting

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
)

func expandAzureRmMonitorAutoScaleSettingProfile(input []any) ([]autoscalesettings.AutoscaleProfile, error) {
	results := make([]autoscalesettings.AutoscaleProfile, 0)

	for _, v := range input {
		raw := v.(map[string]any)

		name := raw["name"].(string)

		// this is Required, so we don't need to check for optionals here
		capacitiesRaw := raw["capacity"].([]any)
		capacityRaw := capacitiesRaw[0].(map[string]any)
		capacity := autoscalesettings.ScaleCapacity{
			Minimum: strconv.Itoa(capacityRaw["minimum"].(int)),
			Maximum: strconv.Itoa(capacityRaw["maximum"].(int)),
			Default: strconv.Itoa(capacityRaw["default"].(int)),
		}

		recurrencesRaw := raw["recurrence"].([]any)
		recurrence := expandAzureRmMonitorAutoScaleSettingRecurrence(recurrencesRaw)

		rulesRaw := raw["rule"].([]any)
		rules := expandAzureRmMonitorAutoScaleSettingRule(rulesRaw)

		fixedDatesRaw := raw["fixed_date"].([]any)
		fixedDate, err := expandAzureRmMonitorAutoScaleSettingFixedDate(fixedDatesRaw)
		if err != nil {
			return nil, fmt.Errorf("expanding `fixed_date`: %+v", err)
		}

		result := autoscalesettings.AutoscaleProfile{
			Name:       name,
			Capacity:   capacity,
			FixedDate:  fixedDate,
			Recurrence: recurrence,
			Rules:      rules,
		}
		results = append(results, result)
	}

	return results, nil
}

func expandAzureRmMonitorAutoScaleSettingPredictive(input []any) *autoscalesettings.PredictiveAutoscalePolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	predictive := autoscalesettings.PredictiveAutoscalePolicy{
		ScaleMode: autoscalesettings.PredictiveAutoscalePolicyScaleMode(raw["scale_mode"].(string)),
	}

	if lookAheadTime := raw["look_ahead_time"].(string); lookAheadTime != "" {
		predictive.ScaleLookAheadTime = pointer.To(lookAheadTime)
	}

	return &predictive
}

func expandAzureRmMonitorAutoScaleSettingRule(input []any) []autoscalesettings.ScaleRule {
	rules := make([]autoscalesettings.ScaleRule, 0)

	for _, v := range input {
		ruleRaw := v.(map[string]any)

		triggersRaw := ruleRaw["metric_trigger"].([]any)
		triggerRaw := triggersRaw[0].(map[string]any)
		metricTrigger := autoscalesettings.MetricTrigger{
			MetricName:        triggerRaw["metric_name"].(string),
			MetricNamespace:   pointer.To(triggerRaw["metric_namespace"].(string)),
			MetricResourceUri: triggerRaw["metric_resource_id"].(string),
			TimeGrain:         triggerRaw["time_grain"].(string),
			Statistic:         autoscalesettings.MetricStatisticType(triggerRaw["statistic"].(string)),
			TimeWindow:        triggerRaw["time_window"].(string),
			TimeAggregation:   autoscalesettings.TimeAggregationType(triggerRaw["time_aggregation"].(string)),
			Operator:          autoscalesettings.ComparisonOperationType(triggerRaw["operator"].(string)),
			Threshold:         triggerRaw["threshold"].(float64),
			Dimensions:        expandAzureRmMonitorAutoScaleSettingRuleDimensions(triggerRaw["dimensions"].([]any)),
			DividePerInstance: pointer.To(triggerRaw["divide_by_instance_count"].(bool)),
		}

		actionsRaw := ruleRaw["scale_action"].([]any)
		actionRaw := actionsRaw[0].(map[string]any)
		scaleAction := autoscalesettings.ScaleAction{
			Direction: autoscalesettings.ScaleDirection(actionRaw["direction"].(string)),
			Type:      autoscalesettings.ScaleType(actionRaw["type"].(string)),
			Value:     pointer.To(strconv.Itoa(actionRaw["value"].(int))),
			Cooldown:  actionRaw["cooldown"].(string),
		}

		rule := autoscalesettings.ScaleRule{
			MetricTrigger: metricTrigger,
			ScaleAction:   scaleAction,
		}

		rules = append(rules, rule)
	}

	return rules
}

func expandAzureRmMonitorAutoScaleSettingFixedDate(input []any) (*autoscalesettings.TimeWindow, error) {
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]any)

	startString := raw["start"].(string)
	startTime, err := date.ParseTime(time.RFC3339, startString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse `start` time %q as an RFC3339 date: %+v", startString, err)
	}
	endString := raw["end"].(string)
	endTime, err := date.ParseTime(time.RFC3339, endString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse `end` time %q as an RFC3339 date: %+v", endString, err)
	}

	timeZone := raw["timezone"].(string)
	timeWindow := autoscalesettings.TimeWindow{
		TimeZone: pointer.To(timeZone),
	}

	timeWindow.SetStartAsTime(startTime)
	timeWindow.SetEndAsTime(endTime)

	return &timeWindow, nil
}

func expandAzureRmMonitorAutoScaleSettingRecurrence(input []any) *autoscalesettings.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrenceRaw := input[0].(map[string]any)

	timeZone := recurrenceRaw["timezone"].(string)
	days := make([]string, 0)
	for _, dayItem := range recurrenceRaw["days"].([]any) {
		days = append(days, dayItem.(string))
	}

	hours := make([]int64, 0)
	for _, hourItem := range recurrenceRaw["hours"].([]any) {
		hours = append(hours, int64(hourItem.(int)))
	}

	minutes := make([]int64, 0)
	for _, minuteItem := range recurrenceRaw["minutes"].([]any) {
		minutes = append(minutes, int64(minuteItem.(int)))
	}

	return &autoscalesettings.Recurrence{
		// API docs say this has to be `Week`.
		Frequency: autoscalesettings.RecurrenceFrequencyWeek,
		Schedule: autoscalesettings.RecurrentSchedule{
			TimeZone: timeZone,
			Days:     days,
			Hours:    hours,
			Minutes:  minutes,
		},
	}
}

func expandAzureRmMonitorAutoScaleSettingNotifications(input []any) *[]autoscalesettings.AutoscaleNotification {
	notifications := make([]autoscalesettings.AutoscaleNotification, 0)

	for _, v := range input {
		notificationRaw := v.(map[string]any)

		configsRaw := notificationRaw["webhook"].([]any)
		webhooks := expandAzureRmMonitorAutoScaleSettingNotificationWebhook(configsRaw)

		notification := autoscalesettings.AutoscaleNotification{
			Operation: "scale",
			WebHooks:  webhooks,
		}

		emailsRaw := notificationRaw["email"].([]any)
		if len(emailsRaw) > 0 && emailsRaw[0] != nil {
			notification.Email = expandAzureRmMonitorAutoScaleSettingNotificationEmail(emailsRaw[0].(map[string]any))
		}

		notifications = append(notifications, notification)
	}

	return &notifications
}

func expandAzureRmMonitorAutoScaleSettingNotificationEmail(input map[string]any) *autoscalesettings.EmailNotification {
	customEmails := make([]string, 0)
	if v, ok := input["custom_emails"]; ok {
		for _, item := range v.([]any) {
			customEmails = append(customEmails, item.(string))
		}
	}

	email := autoscalesettings.EmailNotification{
		CustomEmails:                       &customEmails,
		SendToSubscriptionAdministrator:    pointer.To(input["send_to_subscription_administrator"].(bool)),
		SendToSubscriptionCoAdministrators: pointer.To(input["send_to_subscription_co_administrator"].(bool)),
	}

	return &email
}

func expandAzureRmMonitorAutoScaleSettingNotificationWebhook(input []any) *[]autoscalesettings.WebhookNotification {
	webhooks := make([]autoscalesettings.WebhookNotification, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		webhookRaw := v.(map[string]any)

		webhook := autoscalesettings.WebhookNotification{
			ServiceUri: pointer.To(webhookRaw["service_uri"].(string)),
		}

		if props, ok := webhookRaw["properties"]; ok {
			properties := make(map[string]string)
			for key, value := range props.(map[string]any) {
				properties[key] = value.(string)
			}

			webhook.Properties = &properties
		}

		webhooks = append(webhooks, webhook)
	}

	return &webhooks
}

func expandAzureRmMonitorAutoScaleSettingRuleDimensions(input []any) *[]autoscalesettings.ScaleRuleMetricDimension {
	dimensions := make([]autoscalesettings.ScaleRuleMetricDimension, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		dimensionRaw := v.(map[string]any)

		dimension := autoscalesettings.ScaleRuleMetricDimension{
			DimensionName: dimensionRaw["name"].(string),
			Operator:      autoscalesettings.ScaleRuleMetricDimensionOperationType(dimensionRaw["operator"].(string)),
			Values:        helpers.ExpandStringValues(dimensionRaw["values"].([]any)),
		}

		dimensions = append(dimensions, dimension)
	}

	return &dimensions
}
