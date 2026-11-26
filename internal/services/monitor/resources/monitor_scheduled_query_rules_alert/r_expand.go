// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandMonitorScheduledQueryRulesAlertingAction(d *pluginsdk.ResourceData) *scheduledqueryrules.AlertingAction {
	alertActionRaw := d.Get("action").([]any)
	alertAction := expandMonitorScheduledQueryRulesAlertAction(alertActionRaw)
	severityRaw := d.Get("severity").(int)
	severity := strconv.Itoa(severityRaw)

	triggerRaw := d.Get("trigger").([]any)
	trigger := expandMonitorScheduledQueryRulesAlertTrigger(triggerRaw)

	action := scheduledqueryrules.AlertingAction{
		AznsAction: alertAction,
		Severity:   scheduledqueryrules.AlertSeverity(severity),
		Trigger:    trigger,
	}

	if throttling, ok := d.Get("throttling").(int); ok && throttling != 0 {
		action.ThrottlingInMin = pointer.To(int64(throttling))
	}

	return &action
}

func expandMonitorScheduledQueryRulesAlertAction(input []any) *scheduledqueryrules.AzNsActionGroup {
	result := scheduledqueryrules.AzNsActionGroup{}

	if len(input) == 0 {
		return &result
	}
	for _, item := range input {
		if item == nil {
			continue
		}

		v, ok := item.(map[string]any)
		if !ok {
			continue
		}
		actionGroups := v["action_group"].(*pluginsdk.Set).List()
		result.ActionGroup = pluginsdk.ExpandStringSlice(actionGroups)
		result.EmailSubject = pointer.To(v["email_subject"].(string))
		if v := v["custom_webhook_payload"].(string); v != "" {
			result.CustomWebhookPayload = pointer.To(v)
		}
	}

	return &result
}

func expandMonitorScheduledQueryRulesAlertMetricTrigger(input []any) *scheduledqueryrules.LogMetricTrigger {
	if len(input) == 0 {
		return nil
	}

	result := scheduledqueryrules.LogMetricTrigger{}
	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]any)
		if !ok {
			continue
		}
		result.ThresholdOperator = pointer.ToEnum[scheduledqueryrules.ConditionalOperator](v["operator"].(string))
		result.Threshold = pointer.To(v["threshold"].(float64))
		result.MetricTriggerType = pointer.ToEnum[scheduledqueryrules.MetricTriggerType](v["metric_trigger_type"].(string))
		result.MetricColumn = pointer.To(v["metric_column"].(string))
	}

	return &result
}

func expandMonitorScheduledQueryRulesAlertSchedule(d *pluginsdk.ResourceData) *scheduledqueryrules.Schedule {
	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)

	schedule := scheduledqueryrules.Schedule{
		FrequencyInMinutes:  int64(frequency),
		TimeWindowInMinutes: int64(timeWindow),
	}

	return &schedule
}

func expandMonitorScheduledQueryRulesAlertTrigger(input []any) scheduledqueryrules.TriggerCondition {
	result := scheduledqueryrules.TriggerCondition{}
	if len(input) == 0 {
		return result
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metricTriggerRaw := v["metric_trigger"].([]any)

		result.ThresholdOperator = scheduledqueryrules.ConditionalOperator(v["operator"].(string))
		result.Threshold = v["threshold"].(float64)
		result.MetricTrigger = expandMonitorScheduledQueryRulesAlertMetricTrigger(metricTriggerRaw)
	}

	return result
}
