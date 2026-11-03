// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_global_vm_shutdown_schedule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestGlobalVMShutdownScheduleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmID := d.Get("virtual_machine_id").(string)
	vmId, err := commonids.ParseVirtualMachineIDInsensitively(vmID)
	if err != nil {
		return err
	}

	// Can't find any official documentation on this, but the API returns a 400 for any other name.
	// The best example I could find is here: https://social.msdn.microsoft.com/Forums/en-US/25a02403-dba9-4bcb-bdcc-1f4afcba5b65/powershell-script-to-autoshutdown-azure-virtual-machine?forum=WAVirtualMachinesforWindows
	name := "shutdown-computevm-" + vmId.VirtualMachineName
	id := globalschedules.NewScheduleID(vmId.SubscriptionId, vmId.ResourceGroupName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, globalschedules.GetOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_dev_test_global_vm_shutdown_schedule", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))

	schedule := globalschedules.Schedule{
		Location: &location,
		Properties: globalschedules.ScheduleProperties{
			TargetResourceId: &vmID,
			TaskType:         pointer.To("ComputeVmShutdownTask"),
		},
		Tags: helpers.ExpandTags(d.Get("tags").(map[string]any)),
	}

	statusEnabled := globalschedules.EnableStatusDisabled
	if d.Get("enabled").(bool) {
		statusEnabled = globalschedules.EnableStatusEnabled
	}
	schedule.Properties.Status = &statusEnabled

	if timeZoneId := d.Get("timezone").(string); timeZoneId != "" {
		schedule.Properties.TimeZoneId = &timeZoneId
	}

	if v, ok := d.GetOk("daily_recurrence_time"); ok {
		schedule.Properties.DailyRecurrence = expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(v)
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		schedule.Properties.NotificationSettings = expandDevTestGlobalVMShutdownScheduleNotificationSettings(d)
	}

	if _, err := client.CreateOrUpdate(ctx, id, schedule); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceDevTestGlobalVMShutdownScheduleRead(d, meta)
}

func expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyTime any) *globalschedules.DayDetails {
	return &globalschedules.DayDetails{
		Time: pointer.To(dailyTime.(string)),
	}
}

func expandDevTestGlobalVMShutdownScheduleNotificationSettings(d *pluginsdk.ResourceData) *globalschedules.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]any)
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]any)

	var notificationStatus globalschedules.EnableStatus
	if notificationSettingsConfig["enabled"].(bool) {
		notificationStatus = globalschedules.EnableStatusEnabled
	} else {
		notificationStatus = globalschedules.EnableStatusDisabled
	}

	return &globalschedules.NotificationSettings{
		WebhookURL:     pointer.To(notificationSettingsConfig["webhook_url"].(string)),
		TimeInMinutes:  pointer.To(int64(notificationSettingsConfig["time_in_minutes"].(int))),
		Status:         &notificationStatus,
		EmailRecipient: pointer.To(notificationSettingsConfig["email"].(string)),
	}
}
