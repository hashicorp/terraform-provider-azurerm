// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestLabSchedulesCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schedules.NewLabScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, schedules.GetOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_dev_test_schedule", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))

	schedule := schedules.Schedule{
		Location:   &location,
		Properties: schedules.ScheduleProperties{},
		Tags:       helpers.ExpandTags(d.Get("tags").(map[string]any)),
	}

	status := schedules.EnableStatusDisabled
	if d.Get("status").(string) == string(schedules.EnableStatusEnabled) {
		status = schedules.EnableStatusEnabled
	}
	schedule.Properties.Status = &status

	if taskType := d.Get("task_type").(string); taskType != "" {
		schedule.Properties.TaskType = &taskType
	}

	if timeZoneId := d.Get("time_zone_id").(string); timeZoneId != "" {
		schedule.Properties.TimeZoneId = &timeZoneId
	}

	if v, ok := d.GetOk("weekly_recurrence"); ok {
		schedule.Properties.WeeklyRecurrence = expandDevTestScheduleRecurrenceWeekly(v)
	}

	if v, ok := d.GetOk("daily_recurrence"); ok {
		schedule.Properties.DailyRecurrence = expandDevTestScheduleRecurrenceDaily(v)
	}

	if v, ok := d.GetOk("hourly_recurrence"); ok {
		schedule.Properties.HourlyRecurrence = expandDevTestScheduleRecurrenceHourly(v)
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		schedule.Properties.NotificationSettings = expandDevTestScheduleNotificationSettings(d)
	}

	if _, err := client.CreateOrUpdate(ctx, id, schedule); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceDevTestLabSchedulesRead(d, meta)
}
