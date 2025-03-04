// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestGlobalVMShutdownSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Read:   resourceDevTestGlobalVMShutdownScheduleRead,
		Update: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Delete: resourceDevTestGlobalVMShutdownScheduleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := globalschedules.ParseScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualMachineID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"daily_recurrence_time": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
					"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
				),
			},

			"timezone": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: computeValidate.VirtualMachineTimeZoneCaseInsensitive(),
			},

			"notification_settings": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"time_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(15, 120),
						},
						"webhook_url": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"email": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDevTestGlobalVMShutdownScheduleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmID := d.Get("virtual_machine_id").(string)
	vmId, err := commonids.ParseVirtualMachineID(vmID)
	if err != nil {
		return err
	}

	// Can't find any official documentation on this, but the API returns a 400 for any other name.
	// The best example I could find is here: https://social.msdn.microsoft.com/Forums/en-US/25a02403-dba9-4bcb-bdcc-1f4afcba5b65/powershell-script-to-autoshutdown-azure-virtual-machine?forum=WAVirtualMachinesforWindows
	name := "shutdown-computevm-" + vmId.VirtualMachineName
	id := globalschedules.NewScheduleID(vmId.SubscriptionId, vmId.ResourceGroupName, name)

	if d.IsNewResource() {
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

	location := azure.NormalizeLocation(d.Get("location").(string))
	taskType := "ComputeVmShutdownTask"

	schedule := globalschedules.Schedule{
		Location: &location,
		Properties: globalschedules.ScheduleProperties{
			TargetResourceId: &vmID,
			TaskType:         &taskType,
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
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
		dailyRecurrence := expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(v)
		schedule.Properties.DailyRecurrence = dailyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings := expandDevTestGlobalVMShutdownScheduleNotificationSettings(d)
		schedule.Properties.NotificationSettings = notificationSettings
	}

	if _, err := client.CreateOrUpdate(ctx, id, schedule); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceDevTestGlobalVMShutdownScheduleRead(d, meta)
}

func resourceDevTestGlobalVMShutdownScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if location := resp.Model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
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
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceDevTestGlobalVMShutdownScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := globalschedules.ParseScheduleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return err
	}

	return nil
}

func expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyTime interface{}) *globalschedules.DayDetails {
	time := dailyTime.(string)
	return &globalschedules.DayDetails{
		Time: &time,
	}
}

func flattenDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyRecurrence *globalschedules.DayDetails) interface{} {
	if dailyRecurrence == nil {
		return nil
	}

	var result string
	if dailyRecurrence.Time != nil {
		result = *dailyRecurrence.Time
	}

	return result
}

func expandDevTestGlobalVMShutdownScheduleNotificationSettings(d *pluginsdk.ResourceData) *globalschedules.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]interface{})
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]interface{})
	webhookURL := notificationSettingsConfig["webhook_url"].(string)
	timeInMinutes := int64(notificationSettingsConfig["time_in_minutes"].(int))
	email := notificationSettingsConfig["email"].(string)

	var notificationStatus globalschedules.EnableStatus
	if notificationSettingsConfig["enabled"].(bool) {
		notificationStatus = globalschedules.EnableStatusEnabled
	} else {
		notificationStatus = globalschedules.EnableStatusDisabled
	}

	return &globalschedules.NotificationSettings{
		WebhookURL:     &webhookURL,
		TimeInMinutes:  &timeInMinutes,
		Status:         &notificationStatus,
		EmailRecipient: &email,
	}
}

func flattenDevTestGlobalVMShutdownScheduleNotificationSettings(notificationSettings *globalschedules.NotificationSettings) []interface{} {
	if notificationSettings == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if notificationSettings.WebhookURL != nil {
		result["webhook_url"] = *notificationSettings.WebhookURL
	}

	if notificationSettings.TimeInMinutes != nil {
		result["time_in_minutes"] = *notificationSettings.TimeInMinutes
	}

	result["enabled"] = *notificationSettings.Status == globalschedules.EnableStatusEnabled
	result["email"] = notificationSettings.EmailRecipient

	return []interface{}{result}
}
