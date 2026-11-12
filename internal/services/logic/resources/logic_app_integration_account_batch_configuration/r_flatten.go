// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_batch_configuration

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenIntegrationAccountBatchConfigurationBatchReleaseCriteria(input integrationaccountbatchconfigurations.BatchReleaseCriteria) []any {
	return []any{
		map[string]any{
			"batch_size":    pointer.From(input.BatchSize),
			"message_count": pointer.From(input.MessageCount),
			"recurrence":    flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input.Recurrence),
		},
	}
}

func flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input *integrationaccountbatchconfigurations.WorkflowTriggerRecurrence) []any {
	if input == nil {
		return make([]any, 0)
	}

	endTime := pointer.From(input.EndTime)

	var frequency integrationaccountbatchconfigurations.RecurrenceFrequency
	if input.Frequency != nil && *input.Frequency != "" {
		frequency = *input.Frequency
	}

	interval := pointer.From(input.Interval)

	startTime := pointer.From(input.StartTime)

	timeZone := pointer.From(input.TimeZone)

	return []any{
		map[string]any{
			"end_time":   endTime,
			"frequency":  frequency,
			"interval":   interval,
			"schedule":   flattenIntegrationAccountBatchConfigurationRecurrenceSchedule(input.Schedule),
			"start_time": startTime,
			"time_zone":  timeZone,
		},
	}
}

func flattenIntegrationAccountBatchConfigurationRecurrenceSchedule(input *integrationaccountbatchconfigurations.RecurrenceSchedule) []any {
	if input == nil {
		return make([]any, 0)
	}

	var weekDays []any
	if input.WeekDays != nil {
		weekDaysCast := make([]string, 0)
		for _, item := range *input.WeekDays {
			weekDaysCast = append(weekDaysCast, string(item))
		}
		weekDays = pluginsdk.FlattenSlice(&weekDaysCast)
	}

	return []any{
		map[string]any{
			"hours":      pluginsdk.FlattenSlice(input.Hours),
			"minutes":    pluginsdk.FlattenSlice(input.Minutes),
			"month_days": pluginsdk.FlattenSlice(input.MonthDays),
			"monthly":    flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input.MonthlyOccurrences),
			"week_days":  weekDays,
		},
	}
}

func flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input *[]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var day integrationaccountbatchconfigurations.DayOfWeek
		if item.Day != nil && *item.Day != "" {
			day = *item.Day
		}

		occurrence := pointer.From(item.Occurrence)

		results = append(results, map[string]any{
			"weekday": day,
			"week":    occurrence,
		})
	}

	return results
}
