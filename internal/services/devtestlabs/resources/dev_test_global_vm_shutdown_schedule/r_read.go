// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_global_vm_shutdown_schedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestGlobalVMShutdownScheduleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := globalschedules.ParseScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, globalschedules.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %s", *id, err)
	}

	if model := resp.Model; model != nil {
		if loc := resp.Model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}

		props := resp.Model.Properties
		d.Set("virtual_machine_id", props.TargetResourceId)
		d.Set("timezone", props.TimeZoneId)
		d.Set("enabled", *props.Status == globalschedules.EnableStatusEnabled)

		if err := d.Set("daily_recurrence_time", flattenDevTestGlobalVMShutdownScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
			return fmt.Errorf("setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("notification_settings", flattenDevTestGlobalVMShutdownScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("setting `notificationSettings`: %#v", err)
		}
		if err = tags.FlattenAndSet(d, helpers.FlattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func flattenDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyRecurrence *globalschedules.DayDetails) any {
	if dailyRecurrence == nil {
		return nil
	}

	return pointer.From(dailyRecurrence.Time)
}

func flattenDevTestGlobalVMShutdownScheduleNotificationSettings(notificationSettings *globalschedules.NotificationSettings) []any {
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

	result["enabled"] = *notificationSettings.Status == globalschedules.EnableStatusEnabled
	result["email"] = notificationSettings.EmailRecipient

	return []any{result}
}
