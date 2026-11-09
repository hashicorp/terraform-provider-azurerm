// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
)

func flattenAzureRmDevTestLabScheduleRecurrenceDaily(dailyRecurrence *schedules.DayDetails) []any {
	if dailyRecurrence == nil {
		return []any{}
	}

	result := make(map[string]any)

	if dailyRecurrence.Time != nil {
		result["time"] = *dailyRecurrence.Time
	}

	return []any{result}
}

func flattenAzureRmDevTestLabScheduleRecurrenceWeekly(weeklyRecurrence *schedules.WeekDetails) []any {
	if weeklyRecurrence == nil {
		return []any{}
	}

	result := make(map[string]any)

	if weeklyRecurrence.Time != nil {
		result["time"] = *weeklyRecurrence.Time
	}

	weekDays := make([]string, 0)
	if w := weeklyRecurrence.Weekdays; w != nil {
		weekDays = *w
	}
	result["week_days"] = weekDays

	return []any{result}
}

func flattenAzureRmDevTestLabScheduleRecurrenceHourly(hourlyRecurrence *schedules.HourDetails) []any {
	if hourlyRecurrence == nil {
		return []any{}
	}

	result := make(map[string]any)

	if hourlyRecurrence.Minute != nil {
		result["minute"] = *hourlyRecurrence.Minute
	}

	return []any{result}
}

func flattenAzureRmDevTestLabScheduleNotificationSettings(notificationSettings *schedules.NotificationSettings) []any {
	if notificationSettings == nil {
		return []any{}
	}

	result := make(map[string]any)

	if notificationSettings.WebhookURL != nil {
		result["webhook_url"] = *notificationSettings.WebhookURL
	}

	if notificationSettings.TimeInMinutes != nil {
		result["time_in_minutes"] = *notificationSettings.TimeInMinutes
	}

	if notificationSettings.Status != nil {
		result["status"] = *notificationSettings.Status
	}

	return []any{result}
}
