package devtestlabs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevTestLabGlobalVMShutdownSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevTestLabGlobalVMShutdownScheduleCreateUpdate,
		Read:   resourceArmDevTestLabGlobalVMShutdownScheduleRead,
		Update: resourceArmDevTestLabGlobalVMShutdownScheduleCreateUpdate,
		Delete: resourceArmDevTestLabGlobalVMShutdownScheduleDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ScheduleID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"virtual_machine_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"daily_recurrence_time": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
					"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
				),
			},

			"timezone": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.VirtualMachineTimeZoneCaseInsensitive(),
			},

			"notification_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"time_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(15, 120),
						},
						"webhook_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDevTestLabGlobalVMShutdownScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		dailyRecurrence := expandArmDevTestLabGlobalVMShutdownScheduleRecurrenceDaily(v)
		schedule.DailyRecurrence = dailyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings := expandArmDevTestLabGlobalVMShutdownScheduleNotificationSettings(d)
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

	return resourceArmDevTestLabGlobalVMShutdownScheduleRead(d, meta)
}

func resourceArmDevTestLabGlobalVMShutdownScheduleRead(d *schema.ResourceData, meta interface{}) error {
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

		if err := d.Set("daily_recurrence_time", flattenArmDevTestLabGlobalVMShutdownScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
			return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("notification_settings", flattenArmDevTestLabGlobalVMShutdownScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notificationSettings`: %#v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDevTestLabGlobalVMShutdownScheduleDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandArmDevTestLabGlobalVMShutdownScheduleRecurrenceDaily(dailyTime interface{}) *dtl.DayDetails {
	time := dailyTime.(string)
	return &dtl.DayDetails{
		Time: &time,
	}
}

func flattenArmDevTestLabGlobalVMShutdownScheduleRecurrenceDaily(dailyRecurrence *dtl.DayDetails) interface{} {
	if dailyRecurrence == nil {
		return nil
	}

	var result string
	if dailyRecurrence.Time != nil {
		result = *dailyRecurrence.Time
	}

	return result
}

func expandArmDevTestLabGlobalVMShutdownScheduleNotificationSettings(d *schema.ResourceData) *dtl.NotificationSettings {
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

func flattenArmDevTestLabGlobalVMShutdownScheduleNotificationSettings(notificationSettings *dtl.NotificationSettings) []interface{} {
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
