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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDevTestSchedules() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevTestSchedulesCreateUpdate,
		Read:   resourceDevTestSchedulesRead,
		Update: resourceDevTestSchedulesCreateUpdate,
		Delete: resourceDevTestSchedulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"lab_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  dtl.EnableStatusDisabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.EnableStatusEnabled),
					string(dtl.EnableStatusDisabled),
				}, false),
			},

			"task_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"weekly_recurrence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
								"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
							),
						},

						"week_days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Monday",
									"Tuesday",
									"Wednesday",
									"Thursday",
									"Friday",
									"Saturday",
									"Sunday",
								}, false),
							},
						},
					},
				},
			},

			"daily_recurrence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
								"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
							),
						},
					},
				},
			},

			"hourly_recurrence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minute": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"time_zone_id": {
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
						"status": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  dtl.NotificationStatusDisabled,
							ValidateFunc: validation.StringInSlice([]string{
								string(dtl.NotificationStatusEnabled),
								string(dtl.NotificationStatusDisabled),
							}, false),
						},
						"time_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
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

func resourceDevTestSchedulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	devTestLabName := d.Get("lab_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, devTestLabName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Schedule %q (Dev Test Lab %q / Resource Group %q): %s", name, devTestLabName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dev_test_schedule", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	schedule := dtl.Schedule{
		Location:           &location,
		ScheduleProperties: &dtl.ScheduleProperties{},
		Tags:               tags.Expand(t),
	}

	switch status := d.Get("status"); status {
	case string(dtl.EnableStatusEnabled):
		schedule.ScheduleProperties.Status = dtl.EnableStatusEnabled
	case string(dtl.EnableStatusDisabled):
		schedule.ScheduleProperties.Status = dtl.EnableStatusDisabled
	default:
	}

	if taskType := d.Get("task_type").(string); taskType != "" {
		schedule.ScheduleProperties.TaskType = &taskType
	}

	if timeZoneId := d.Get("time_zone_id").(string); timeZoneId != "" {
		schedule.ScheduleProperties.TimeZoneID = &timeZoneId
	}

	if v, ok := d.GetOk("weekly_recurrence"); ok {
		weekRecurrence := expandDevTestScheduleRecurrenceWeekly(v)

		schedule.WeeklyRecurrence = weekRecurrence
	}

	if v, ok := d.GetOk("daily_recurrence"); ok {
		dailyRecurrence := expandDevTestScheduleRecurrenceDaily(v)
		schedule.DailyRecurrence = dailyRecurrence
	}

	if v, ok := d.GetOk("hourly_recurrence"); ok {
		hourlyRecurrence := expandDevTestScheduleRecurrenceHourly(v)

		schedule.HourlyRecurrence = hourlyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings := expandDevTestScheduleNotificationSettings(d)
		schedule.NotificationSettings = notificationSettings
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, devTestLabName, name, schedule); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, devTestLabName, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read  Dev Test Lab Schedule %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceDevTestSchedulesRead(d, meta)
}

func resourceDevTestSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	devTestLabName := id.Path["labs"]
	name := id.Path["schedules"]

	resp, err := client.Get(ctx, resGroup, devTestLabName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Dev Test Lab Schedule %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("lab_name", devTestLabName)
	d.Set("resource_group_name", resGroup)
	d.Set("task_type", resp.TaskType)

	if props := resp.ScheduleProperties; props != nil {
		d.Set("time_zone_id", props.TimeZoneID)

		d.Set("status", string(props.Status))

		if err := d.Set("weekly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceWeekly(props.WeeklyRecurrence)); err != nil {
			return fmt.Errorf("Error setting `weeklyRecurrence`: %#v", err)
		}

		if err := d.Set("daily_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
			return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("hourly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceHourly(props.HourlyRecurrence)); err != nil {
			return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
		}

		if err := d.Set("notification_settings", flattenAzureRmDevTestLabScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notificationSettings`: %#v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDevTestSchedulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMExtensionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["schedules"]
	devTestLabName := id.Path["labs"]

	future, err := client.Delete(ctx, resGroup, devTestLabName, name)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}

func expandDevTestScheduleRecurrenceDaily(recurrence interface{}) *dtl.DayDetails {
	dailyRecurrenceConfigs := recurrence.([]interface{})
	dailyRecurrenceConfig := dailyRecurrenceConfigs[0].(map[string]interface{})
	dailyTime := dailyRecurrenceConfig["time"].(string)

	return &dtl.DayDetails{
		Time: &dailyTime,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceDaily(dailyRecurrence *dtl.DayDetails) []interface{} {
	if dailyRecurrence == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if dailyRecurrence.Time != nil {
		result["time"] = *dailyRecurrence.Time
	}

	return []interface{}{result}
}

func expandDevTestScheduleRecurrenceWeekly(recurrence interface{}) *dtl.WeekDetails {
	weeklyRecurrenceConfigs := recurrence.([]interface{})
	weeklyRecurrenceConfig := weeklyRecurrenceConfigs[0].(map[string]interface{})
	weeklyTime := weeklyRecurrenceConfig["time"].(string)

	weekDays := make([]string, 0)
	for _, dayItem := range weeklyRecurrenceConfig["week_days"].([]interface{}) {
		weekDays = append(weekDays, dayItem.(string))
	}

	return &dtl.WeekDetails{
		Time:     &weeklyTime,
		Weekdays: &weekDays,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceWeekly(weeklyRecurrence *dtl.WeekDetails) []interface{} {
	if weeklyRecurrence == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if weeklyRecurrence.Time != nil {
		result["time"] = *weeklyRecurrence.Time
	}

	weekDays := make([]string, 0)
	if w := weeklyRecurrence.Weekdays; w != nil {
		weekDays = *w
	}
	result["week_days"] = weekDays

	return []interface{}{result}
}

func expandDevTestScheduleRecurrenceHourly(recurrence interface{}) *dtl.HourDetails {
	hourlyRecurrenceConfigs := recurrence.([]interface{})
	hourlyRecurrenceConfig := hourlyRecurrenceConfigs[0].(map[string]interface{})
	hourlyMinute := int32(hourlyRecurrenceConfig["minute"].(int))

	return &dtl.HourDetails{
		Minute: &hourlyMinute,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceHourly(hourlyRecurrence *dtl.HourDetails) []interface{} {
	if hourlyRecurrence == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if hourlyRecurrence.Minute != nil {
		result["minute"] = *hourlyRecurrence.Minute
	}

	return []interface{}{result}
}

func expandDevTestScheduleNotificationSettings(d *schema.ResourceData) *dtl.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]interface{})
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]interface{})
	webhookUrl := notificationSettingsConfig["webhook_url"].(string)
	timeInMinutes := int32(notificationSettingsConfig["time_in_minutes"].(int))

	notificationStatus := dtl.NotificationStatus(notificationSettingsConfig["status"].(string))

	return &dtl.NotificationSettings{
		WebhookURL:    &webhookUrl,
		TimeInMinutes: &timeInMinutes,
		Status:        notificationStatus,
	}
}

func flattenAzureRmDevTestLabScheduleNotificationSettings(notificationSettings *dtl.NotificationSettings) []interface{} {
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

	if string(notificationSettings.Status) != "" {
		result["status"] = string(notificationSettings.Status)
	}

	return []interface{}{result}
}
