// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_schedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestLabSchedulesRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schedules.ParseLabScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, schedules.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request %s: %s", *id, err)
	}

	d.Set("name", id.ScheduleName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if loc := model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}

		props := model.Properties
		d.Set("time_zone_id", props.TimeZoneId)
		d.Set("task_type", props.TaskType)
		d.Set("status", pointer.FromEnum(props.Status))

		if err := d.Set("weekly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceWeekly(props.WeeklyRecurrence)); err != nil {
			return fmt.Errorf("setting `weeklyRecurrence`: %#v", err)
		}

		if err := d.Set("daily_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
			return fmt.Errorf("setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("hourly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceHourly(props.HourlyRecurrence)); err != nil {
			return fmt.Errorf("setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("notification_settings", flattenAzureRmDevTestLabScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("setting `notificationSettings`: %#v", err)
		}
		if err = tags.FlattenAndSet(d, helpers.FlattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}
