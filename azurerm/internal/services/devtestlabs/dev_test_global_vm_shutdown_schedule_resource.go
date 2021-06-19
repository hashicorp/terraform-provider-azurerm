package devtestlabs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDevTestGlobalVMShutdownSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Read:   resourceDevTestGlobalVMShutdownScheduleRead,
		Update: resourceDevTestGlobalVMShutdownScheduleCreateUpdate,
		Delete: resourceDevTestGlobalVMShutdownScheduleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": azure.SchemaLocation(),

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineID,
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
	id, err := computeParse.VirtualMachineID(vmID)
	if err != nil {
		return err
	}

	// Can't find any official documentation on this, but the API returns a 400 for any other name.
	// The best example I could find is here: https://social.msdn.microsoft.com/Forums/en-US/25a02403-dba9-4bcb-bdcc-1f4afcba5b65/powershell-script-to-autoshutdown-azure-virtual-machine?forum=WAVirtualMachinesforWindows
	name := "shutdown-computevm-" + id.Name

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Schedule %q (Resource Group %q): %s", name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dev_test_global_vm_shutdown_schedule", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	taskType := "ComputeVmShutdownTask"

	schedule := dtl.Schedule{
		Location: &location,
		ScheduleProperties: &dtl.ScheduleProperties{
			TargetResourceID: &vmID,
			TaskType:         &taskType,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.Get("enabled").(bool) {
		schedule.ScheduleProperties.Status = dtl.EnableStatusEnabled
	} else {
		schedule.ScheduleProperties.Status = dtl.EnableStatusDisabled
	}

	if timeZoneId := d.Get("timezone").(string); timeZoneId != "" {
		schedule.ScheduleProperties.TimeZoneID = &timeZoneId
	}

	if v, ok := d.GetOk("daily_recurrence_time"); ok {
		dailyRecurrence := expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(v)
		schedule.DailyRecurrence = dailyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings := expandDevTestGlobalVMShutdownScheduleNotificationSettings(d)
		schedule.NotificationSettings = notificationSettings
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, name, schedule); err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Dev Test Global Schedule %s (resource group %s) ID", name, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceDevTestGlobalVMShutdownScheduleRead(d, meta)
}

func resourceDevTestGlobalVMShutdownScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Dev Test Global Schedule %s: %s", id.Name, err)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ScheduleProperties; props != nil {
		d.Set("virtual_machine_id", props.TargetResourceID)
		d.Set("timezone", props.TimeZoneID)
		d.Set("enabled", props.Status == dtl.EnableStatusEnabled)

		if err := d.Set("daily_recurrence_time", flattenDevTestGlobalVMShutdownScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
			return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("notification_settings", flattenDevTestGlobalVMShutdownScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notificationSettings`: %#v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDevTestGlobalVMShutdownScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.GlobalLabSchedulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScheduleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return err
	}

	return nil
}

func expandDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyTime interface{}) *dtl.DayDetails {
	time := dailyTime.(string)
	return &dtl.DayDetails{
		Time: &time,
	}
}

func flattenDevTestGlobalVMShutdownScheduleRecurrenceDaily(dailyRecurrence *dtl.DayDetails) interface{} {
	if dailyRecurrence == nil {
		return nil
	}

	var result string
	if dailyRecurrence.Time != nil {
		result = *dailyRecurrence.Time
	}

	return result
}

func expandDevTestGlobalVMShutdownScheduleNotificationSettings(d *pluginsdk.ResourceData) *dtl.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]interface{})
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]interface{})
	webhookUrl := notificationSettingsConfig["webhook_url"].(string)
	timeInMinutes := int32(notificationSettingsConfig["time_in_minutes"].(int))

	var notificationStatus dtl.NotificationStatus
	if notificationSettingsConfig["enabled"].(bool) {
		notificationStatus = dtl.NotificationStatusEnabled
	} else {
		notificationStatus = dtl.NotificationStatusDisabled
	}

	return &dtl.NotificationSettings{
		WebhookURL:    &webhookUrl,
		TimeInMinutes: &timeInMinutes,
		Status:        notificationStatus,
	}
}

func flattenDevTestGlobalVMShutdownScheduleNotificationSettings(notificationSettings *dtl.NotificationSettings) []interface{} {
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

	result["enabled"] = notificationSettings.Status == dtl.NotificationStatusEnabled

	return []interface{}{result}
}
