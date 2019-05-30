package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevTestLabSchedules() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevTestLabSchedulesCreateUpdate,
		Read:   resourceArmDevTestLabSchedulesRead,
		Update: resourceArmDevTestLabSchedulesCreateUpdate,
		Delete: resourceArmDevTestLabSchedulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"dev_test_lab_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  dtl.EnableStatusDisabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.EnableStatusEnabled),
					string(dtl.EnableStatusDisabled),
				}, true),
			},

			"task_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"weekly_recurrence": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Required: true,
						},

						"week_days": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(scheduler.Monday),
									string(scheduler.Tuesday),
									string(scheduler.Wednesday),
									string(scheduler.Thursday),
									string(scheduler.Friday),
									string(scheduler.Saturday),
									string(scheduler.Sunday),
								}, true),
								DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							},
							// Set: set.HashStringIgnoreCase,
						},
					},
				},
			},

			"daily_recurrence": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"hourly_recurrence": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minute": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"time_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"notification_settings": {
				Type:     schema.TypeSet,
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
							}, true),
						},
						"time_in_minutes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"webhook_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDevTestLabSchedulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestLabSchedulesClient

	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	devTestLabName := d.Get("dev_test_lab_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, devTestLabName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Schedule %q (Dev Test Lab %q / Resource Group %q): %s", name, devTestLabName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dev_test_lab_schedule", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	schedule := dtl.Schedule{
		Location:           &location,
		ScheduleProperties: &dtl.ScheduleProperties{},
		Tags:               expandTags(tags),
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

	if _, ok := d.GetOk("weekly_recurrence"); ok {
		weekRecurrence, err2 := expandArmDevTestLabScheduleRecurrenceWeekly(d)
		if err2 != nil {
			return err2
		}

		schedule.WeeklyRecurrence = weekRecurrence
	}

	if _, ok := d.GetOk("daily_recurrence"); ok {
		dailyRecurrence, err2 := expandArmDevTestLabScheduleRecurrenceDaily(d)
		if err2 != nil {
			return err2
		}

		schedule.DailyRecurrence = dailyRecurrence
	}

	if _, ok := d.GetOk("hourly_recurrence"); ok {
		hourlyRecurrence, err2 := expandArmDevTestLabScheduleRecurrenceHourly(d)
		if err2 != nil {
			return err2
		}

		schedule.HourlyRecurrence = hourlyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings, err2 := expandArmDevTestLabScheduleNotificationSettings(d)
		if err2 != nil {
			return err2
		}
		schedule.NotificationSettings = notificationSettings
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, devTestLabName, name, schedule)
	if err != nil {
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

	return resourceArmDevTestLabSchedulesRead(d, meta)
}

func resourceArmDevTestLabSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestLabSchedulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("dev_test_lab_name", devTestLabName)
	d.Set("resource_group_name", resGroup)

	if props := resp.ScheduleProperties; props != nil {
		if timeZoneId := props.TimeZoneID; *timeZoneId != "" {
			d.Set("time_zone_id", *timeZoneId)
		}

		switch status := props.Status; status {
		case dtl.EnableStatusEnabled:
			d.Set("status", string(dtl.EnableStatusEnabled))
		case dtl.EnableStatusDisabled:
			d.Set("status", string(dtl.EnableStatusDisabled))
		default:
		}

		if weeklyRecurrence := props.WeeklyRecurrence; weeklyRecurrence != nil {
			if err := d.Set("weekly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceWeekly(props.WeeklyRecurrence)); err != nil {
				return fmt.Errorf("Error setting `weeklyRecurrence`: %#v", err)
			}
		}

		if dailyRecurrence := props.DailyRecurrence; dailyRecurrence != nil {
			if err := d.Set("daily_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceDaily(props.DailyRecurrence)); err != nil {
				return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
			}
		}

		if hourlyRecurrence := props.HourlyRecurrence; hourlyRecurrence != nil {
			if err := d.Set("hourly_recurrence", flattenAzureRmDevTestLabScheduleRecurrenceHourly(props.HourlyRecurrence)); err != nil {
				return fmt.Errorf("Error setting `dailyRecurrence`: %#v", err)
			}
		}

		if err := d.Set("notification_settings", flattenAzureRmDevTestLabScheduleNotificationSettings(props.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notificationSettings`: %#v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDevTestLabSchedulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmExtensionClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

func expandArmDevTestLabScheduleRecurrenceDaily(d *schema.ResourceData) (*dtl.DayDetails, error) {
	dailyRecurrenceConfigs := d.Get("daily_recurrence").(*schema.Set).List()
	dailyRecurrenceConfig := dailyRecurrenceConfigs[0].(map[string]interface{})
	dailyTime := dailyRecurrenceConfig["time"].(string)

	return &dtl.DayDetails{
		Time: &dailyTime,
	}, nil

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

func expandArmDevTestLabScheduleRecurrenceWeekly(d *schema.ResourceData) (*dtl.WeekDetails, error) {
	weeklyRecurrenceConfigs := d.Get("weekly_recurrence").(*schema.Set).List()
	weeklyRecurrenceConfig := weeklyRecurrenceConfigs[0].(map[string]interface{})

	weeklyTime := weeklyRecurrenceConfig["time"].(string)

	weekDays := make([]string, 0)
	for _, dayItem := range weeklyRecurrenceConfig["week_days"].([]interface{}) {
		weekDays = append(weekDays, dayItem.(string))
	}

	return &dtl.WeekDetails{
		Time:     &weeklyTime,
		Weekdays: &weekDays,
	}, nil

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

func expandArmDevTestLabScheduleRecurrenceHourly(d *schema.ResourceData) (*dtl.HourDetails, error) {
	hourlyRecurrenceConfigs := d.Get("hourly_recurrence").(*schema.Set).List()
	hourlyRecurrenceConfig := hourlyRecurrenceConfigs[0].(map[string]interface{})
	hourlyMinute := int32(hourlyRecurrenceConfig["minute"].(int))

	return &dtl.HourDetails{
		Minute: &hourlyMinute,
	}, nil

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

func expandArmDevTestLabScheduleNotificationSettings(d *schema.ResourceData) (*dtl.NotificationSettings, error) {

	notificationSettingsConfigs := d.Get("notification_settings").(*schema.Set).List()
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]interface{})
	webhookUrl := notificationSettingsConfig["webhook_url"].(string)
	timeInMinutes := int32(notificationSettingsConfig["time_in_minutes"].(int))

	var notificationStatus dtl.NotificationStatus
	switch status := notificationSettingsConfig["status"]; status {
	case string(dtl.NotificationStatusEnabled):
		notificationStatus = dtl.NotificationStatusEnabled
	case string(dtl.NotificationStatusDisabled):
		notificationStatus = dtl.NotificationStatusDisabled
	default:
	}

	return &dtl.NotificationSettings{
		WebhookURL:    &webhookUrl,
		TimeInMinutes: &timeInMinutes,
		Status:        notificationStatus,
	}, nil

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

		var notificationStatus string
		switch notificationSettings.Status {
		case dtl.NotificationStatusEnabled:
			notificationStatus = string(dtl.NotificationStatusEnabled)
		case dtl.NotificationStatusDisabled:
			notificationStatus = string(dtl.NotificationStatusDisabled)
		default:
		}
		result["status"] = notificationStatus
	}

	return []interface{}{result}
}
