// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandDevTestScheduleRecurrenceDaily(recurrence any) *schedules.DayDetails {
	dailyRecurrenceConfigs := recurrence.([]any)
	dailyRecurrenceConfig := dailyRecurrenceConfigs[0].(map[string]any)

	return &schedules.DayDetails{
		Time: pointer.To(dailyRecurrenceConfig["time"].(string)),
	}
}

func expandDevTestScheduleRecurrenceWeekly(recurrence any) *schedules.WeekDetails {
	weeklyRecurrenceConfigs := recurrence.([]any)
	weeklyRecurrenceConfig := weeklyRecurrenceConfigs[0].(map[string]any)

	weekDays := make([]string, 0)
	for _, dayItem := range weeklyRecurrenceConfig["week_days"].([]any) {
		weekDays = append(weekDays, dayItem.(string))
	}

	return &schedules.WeekDetails{
		Time:     pointer.To(weeklyRecurrenceConfig["time"].(string)),
		Weekdays: &weekDays,
	}
}

func expandDevTestScheduleRecurrenceHourly(recurrence any) *schedules.HourDetails {
	hourlyRecurrenceConfigs := recurrence.([]any)
	hourlyRecurrenceConfig := hourlyRecurrenceConfigs[0].(map[string]any)

	return &schedules.HourDetails{
		Minute: pointer.To(int64(hourlyRecurrenceConfig["minute"].(int))),
	}
}

func expandDevTestScheduleNotificationSettings(d *pluginsdk.ResourceData) *schedules.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]any)
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]any)

	return &schedules.NotificationSettings{
		WebhookURL:    pointer.To(notificationSettingsConfig["webhook_url"].(string)),
		TimeInMinutes: pointer.To(int64(notificationSettingsConfig["time_in_minutes"].(int))),
		Status:        pointer.ToEnum[schedules.EnableStatus](notificationSettingsConfig["status"].(string)),
	}
}
