// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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

	enabled := notificationSettingsConfig["enabled"].(bool)
	if !features.SixPointOh() {
		switch {
		case !pluginsdk.IsExplicitlyNullInConfigAt(d, "notification_settings.0.status"):
			enabled = notificationSettingsConfig["status"].(string) == string(schedules.EnableStatusEnabled)
		case pluginsdk.IsExplicitlyNullInConfigAt(d, "notification_settings.0.enabled"):
			enabled = false
		}
	}

	return &schedules.NotificationSettings{
		WebhookURL:    pointer.To(notificationSettingsConfig["webhook_url"].(string)),
		TimeInMinutes: pointer.To(int64(notificationSettingsConfig["time_in_minutes"].(int))),
		Status:        expandDevTestScheduleEnableStatus(enabled),
	}
}

// expandDevTestScheduleStatus is whether the schedule is enabled, from `enabled`, or in 5.x from `status` when that is set. Leaving both out means
// disabled.
func expandDevTestScheduleStatus(d *pluginsdk.ResourceData) *schedules.EnableStatus {
	enabled := d.Get("enabled").(bool)
	if !features.SixPointOh() {
		switch {
		case !pluginsdk.IsExplicitlyNullInConfig(d, "status"):
			enabled = d.Get("status").(string) == string(schedules.EnableStatusEnabled)
		case pluginsdk.IsExplicitlyNullInConfig(d, "enabled"):
			enabled = false
		}
	}

	return expandDevTestScheduleEnableStatus(enabled)
}

func expandDevTestScheduleEnableStatus(enabled bool) *schedules.EnableStatus {
	if enabled {
		return pointer.To(schedules.EnableStatusEnabled)
	}

	return pointer.To(schedules.EnableStatusDisabled)
}
