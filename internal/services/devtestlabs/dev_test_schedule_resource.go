// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	devTestValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDevTestLabSchedules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestLabSchedulesCreateUpdate,
		Read:   resourceDevTestLabSchedulesRead,
		Update: resourceDevTestLabSchedulesCreateUpdate,
		Delete: resourceDevTestLabSchedulesDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schedules.ParseLabScheduleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestLabScheduleUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"lab_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: devTestValidate.DevTestLabName(),
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  schedules.EnableStatusDisabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(schedules.EnableStatusEnabled),
					string(schedules.EnableStatusDisabled),
				}, false),
			},

			"task_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"weekly_recurrence": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"time": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9])[0-5][0-9]$"),
								"Time of day must match the format HHmm where HH is 00-23 and mm is 00-59",
							),
						},

						"week_days": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"time": {
							Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"minute": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"time_zone_id": {
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
						"status": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  schedules.EnableStatusDisabled,
							ValidateFunc: validation.StringInSlice([]string{
								string(schedules.EnableStatusEnabled),
								string(schedules.EnableStatusDisabled),
							}, false),
						},
						"time_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
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

func resourceDevTestLabSchedulesCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schedules.NewLabScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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

	location := azure.NormalizeLocation(d.Get("location").(string))

	schedule := schedules.Schedule{
		Location:   &location,
		Properties: schedules.ScheduleProperties{},
		Tags:       expandTags(d.Get("tags").(map[string]interface{})),
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
		weekRecurrence := expandDevTestScheduleRecurrenceWeekly(v)

		schedule.Properties.WeeklyRecurrence = weekRecurrence
	}

	if v, ok := d.GetOk("daily_recurrence"); ok {
		dailyRecurrence := expandDevTestScheduleRecurrenceDaily(v)
		schedule.Properties.DailyRecurrence = dailyRecurrence
	}

	if v, ok := d.GetOk("hourly_recurrence"); ok {
		hourlyRecurrence := expandDevTestScheduleRecurrenceHourly(v)

		schedule.Properties.HourlyRecurrence = hourlyRecurrence
	}

	if _, ok := d.GetOk("notification_settings"); ok {
		notificationSettings := expandDevTestScheduleNotificationSettings(d)
		schedule.Properties.NotificationSettings = notificationSettings
	}

	if _, err := client.CreateOrUpdate(ctx, id, schedule); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceDevTestLabSchedulesRead(d, meta)
}

func resourceDevTestLabSchedulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		props := model.Properties
		d.Set("time_zone_id", props.TimeZoneId)
		d.Set("task_type", props.TaskType)
		d.Set("status", string(pointer.From(props.Status)))

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
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceDevTestLabSchedulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabSchedulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schedules.ParseLabScheduleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return err
	}
	return nil
}

func expandDevTestScheduleRecurrenceDaily(recurrence interface{}) *schedules.DayDetails {
	dailyRecurrenceConfigs := recurrence.([]interface{})
	dailyRecurrenceConfig := dailyRecurrenceConfigs[0].(map[string]interface{})
	dailyTime := dailyRecurrenceConfig["time"].(string)

	return &schedules.DayDetails{
		Time: &dailyTime,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceDaily(dailyRecurrence *schedules.DayDetails) []interface{} {
	if dailyRecurrence == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if dailyRecurrence.Time != nil {
		result["time"] = *dailyRecurrence.Time
	}

	return []interface{}{result}
}

func expandDevTestScheduleRecurrenceWeekly(recurrence interface{}) *schedules.WeekDetails {
	weeklyRecurrenceConfigs := recurrence.([]interface{})
	weeklyRecurrenceConfig := weeklyRecurrenceConfigs[0].(map[string]interface{})
	weeklyTime := weeklyRecurrenceConfig["time"].(string)

	weekDays := make([]string, 0)
	for _, dayItem := range weeklyRecurrenceConfig["week_days"].([]interface{}) {
		weekDays = append(weekDays, dayItem.(string))
	}

	return &schedules.WeekDetails{
		Time:     &weeklyTime,
		Weekdays: &weekDays,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceWeekly(weeklyRecurrence *schedules.WeekDetails) []interface{} {
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

func expandDevTestScheduleRecurrenceHourly(recurrence interface{}) *schedules.HourDetails {
	hourlyRecurrenceConfigs := recurrence.([]interface{})
	hourlyRecurrenceConfig := hourlyRecurrenceConfigs[0].(map[string]interface{})
	hourlyMinute := int64(hourlyRecurrenceConfig["minute"].(int))

	return &schedules.HourDetails{
		Minute: &hourlyMinute,
	}
}

func flattenAzureRmDevTestLabScheduleRecurrenceHourly(hourlyRecurrence *schedules.HourDetails) []interface{} {
	if hourlyRecurrence == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if hourlyRecurrence.Minute != nil {
		result["minute"] = *hourlyRecurrence.Minute
	}

	return []interface{}{result}
}

func expandDevTestScheduleNotificationSettings(d *pluginsdk.ResourceData) *schedules.NotificationSettings {
	notificationSettingsConfigs := d.Get("notification_settings").([]interface{})
	notificationSettingsConfig := notificationSettingsConfigs[0].(map[string]interface{})
	webhookUrl := notificationSettingsConfig["webhook_url"].(string)
	timeInMinutes := int64(notificationSettingsConfig["time_in_minutes"].(int))

	notificationStatus := schedules.EnableStatus(notificationSettingsConfig["status"].(string))

	return &schedules.NotificationSettings{
		WebhookURL:    &webhookUrl,
		TimeInMinutes: &timeInMinutes,
		Status:        &notificationStatus,
	}
}

func flattenAzureRmDevTestLabScheduleNotificationSettings(notificationSettings *schedules.NotificationSettings) []interface{} {
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

	if notificationSettings.Status != nil {
		result["status"] = *notificationSettings.Status
	}

	return []interface{}{result}
}
