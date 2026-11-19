// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_batch_configuration

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandIntegrationAccountBatchConfigurationBatchReleaseCriteria(input []any) integrationaccountbatchconfigurations.BatchReleaseCriteria {
	result := integrationaccountbatchconfigurations.BatchReleaseCriteria{}
	if len(input) == 0 {
		return result
	}
	v := input[0].(map[string]any)

	if batchSize := v["batch_size"].(int); batchSize != 0 {
		result.BatchSize = pointer.To(int64(batchSize))
	}

	if messageCount := v["message_count"].(int); messageCount != 0 {
		result.MessageCount = pointer.To(int64(messageCount))
	}

	if recurrence := v["recurrence"].([]any); len(recurrence) != 0 {
		result.Recurrence = expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(recurrence)
	}

	return result
}

func expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input []any) *integrationaccountbatchconfigurations.WorkflowTriggerRecurrence {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	result := integrationaccountbatchconfigurations.WorkflowTriggerRecurrence{
		Frequency: pointer.ToEnum[integrationaccountbatchconfigurations.RecurrenceFrequency](v["frequency"].(string)),
		Interval:  pointer.To(int64(v["interval"].(int))),
	}

	if startTime := v["start_time"].(string); startTime != "" {
		result.StartTime = pointer.To(startTime)
	}

	if endTime := v["end_time"].(string); endTime != "" {
		result.EndTime = pointer.To(endTime)
	}

	if timeZone := v["time_zone"].(string); timeZone != "" {
		result.TimeZone = pointer.To(timeZone)
	}

	if schedule := v["schedule"].([]any); len(schedule) != 0 {
		result.Schedule = expandIntegrationAccountBatchConfigurationRecurrenceSchedule(schedule)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceSchedule(input []any) *integrationaccountbatchconfigurations.RecurrenceSchedule {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	result := integrationaccountbatchconfigurations.RecurrenceSchedule{}

	if hours := v["hours"].(*pluginsdk.Set).List(); len(hours) != 0 {
		result.Hours = pluginsdk.ExpandInt64Slice(hours)
	}

	if minutes := v["minutes"].(*pluginsdk.Set).List(); len(minutes) != 0 {
		result.Minutes = pluginsdk.ExpandInt64Slice(minutes)
	}

	if rawWeekDays := v["week_days"].(*pluginsdk.Set).List(); len(rawWeekDays) != 0 {
		weekDays := make([]integrationaccountbatchconfigurations.DaysOfWeek, 0)
		for _, item := range *pluginsdk.ExpandStringSlice(rawWeekDays) {
			weekDays = append(weekDays, integrationaccountbatchconfigurations.DaysOfWeek(item))
		}
		result.WeekDays = &weekDays
	}

	if monthDays := v["month_days"].(*pluginsdk.Set).List(); len(monthDays) != 0 {
		result.MonthDays = pluginsdk.ExpandInt64Slice(monthDays)
	}

	if monthlyOccurrence := v["monthly"].(*pluginsdk.Set).List(); len(monthlyOccurrence) != 0 {
		result.MonthlyOccurrences = expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(monthlyOccurrence)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(input []any) *[]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence {
	results := make([]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence{
			Day:        pointer.ToEnum[integrationaccountbatchconfigurations.DayOfWeek](v["weekday"].(string)),
			Occurrence: pointer.To(int64(v["week"].(int))),
		})
	}

	return &results
}
