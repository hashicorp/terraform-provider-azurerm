// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	// import time/tzdata to embed timezone information in the program
	// add this to resolve https://github.com/hashicorp/terraform-provider-azurerm/issues/20690
	_ "time/tzdata"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azvalidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationScheduleCreateUpdate,
		Read:   resourceAutomationScheduleRead,
		Update: resourceAutomationScheduleCreateUpdate,
		Delete: resourceAutomationScheduleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schedule.ParseScheduleID(id)
			return err
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
				ValidateFunc: validate.ScheduleName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"frequency": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(schedule.ScheduleFrequencyDay),
					string(schedule.ScheduleFrequencyHour),
					string(schedule.ScheduleFrequencyMonth),
					string(schedule.ScheduleFrequencyOneTime),
					string(schedule.ScheduleFrequencyWeek),
				}, false),
			},

			// ignored when frequency is `OneTime`
			"interval": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true, // defaults to 1 if frequency is not OneTime
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"start_time": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339MinuteTime,
				ValidateFunc:     validation.IsRFC3339Time,
				// defaults to now + 7 minutes in create function if not set
			},

			"expiry_time": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true, // same as start time when OneTime, ridiculous value when recurring: "9999-12-31T15:59:00-08:00"
				DiffSuppressFunc: suppress.RFC3339MinuteTime,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"timezone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "Etc/UTC",
				ValidateFunc: azvalidate.AzureTimeZoneString(),
			},

			"week_days": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(schedule.ScheduleDayMonday),
						string(schedule.ScheduleDayTuesday),
						string(schedule.ScheduleDayWednesday),
						string(schedule.ScheduleDayThursday),
						string(schedule.ScheduleDayFriday),
						string(schedule.ScheduleDaySaturday),
						string(schedule.ScheduleDaySunday),
					}, false),
				},
				Set:           set.HashStringIgnoreCase,
				ConflictsWith: []string{"month_days", "monthly_occurrence"},
			},

			"month_days": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeInt,
					ValidateFunc: validation.All(
						validation.IntBetween(-1, 31),
						validation.IntNotInSlice([]int{0}),
					),
				},
				Set:           set.HashInt,
				ConflictsWith: []string{"week_days", "monthly_occurrence"},
			},

			"monthly_occurrence": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(schedule.ScheduleDayMonday),
								string(schedule.ScheduleDayTuesday),
								string(schedule.ScheduleDayWednesday),
								string(schedule.ScheduleDayThursday),
								string(schedule.ScheduleDayFriday),
								string(schedule.ScheduleDaySaturday),
								string(schedule.ScheduleDaySunday),
							}, false),
						},
						"occurrence": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ValidateFunc: validation.All(
								validation.IntBetween(-1, 5),
								validation.IntNotInSlice([]int{0}),
							),
						},
					},
				},
				ConflictsWith: []string{"week_days", "month_days"},
			},
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			frequency := strings.ToLower(diff.Get("frequency").(string))
			interval, _ := diff.GetOk("interval")
			if frequency == "onetime" && interval.(int) > 0 {
				// because `interval` is optional and computed, so interval value can exist even it removed from configuration
				// have to check it in raw config
				intervalVal := diff.GetRawConfig().GetAttr("interval")
				if !intervalVal.IsNull() {
					return fmt.Errorf("`interval` cannot be set when frequency is `OneTime`")
				}
			}

			_, hasWeekDays := diff.GetOk("week_days")
			if hasWeekDays && frequency != "week" {
				return fmt.Errorf("`week_days` can only be set when frequency is `Week`")
			}

			_, hasMonthDays := diff.GetOk("month_days")
			if hasMonthDays && frequency != "month" {
				return fmt.Errorf("`month_days` can only be set when frequency is `Month`")
			}

			_, hasMonthlyOccurrences := diff.GetOk("monthly_occurrence")
			if hasMonthlyOccurrences && frequency != "month" {
				return fmt.Errorf("`monthly_occurrence` can only be set when frequency is `Month`")
			}

			return nil
		}),
	}
}

func resourceAutomationScheduleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Schedule
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Schedule creation.")

	id := schedule.NewScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_schedule", id.ID())
		}
	}

	frequency := d.Get("frequency").(string)
	timeZone := d.Get("timezone").(string)
	description := d.Get("description").(string)

	parameters := schedule.ScheduleCreateOrUpdateParameters{
		Name: id.ScheduleName,
		Properties: schedule.ScheduleCreateOrUpdateProperties{
			Description: &description,
			Frequency:   schedule.ScheduleFrequency(frequency),
			TimeZone:    &timeZone,
		},
	}

	// start time can default to now + 7 (5 could be invalid by the time the API is called)
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}
	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		duration := time.Duration(5) * time.Minute
		if time.Until(t) < duration {
			return fmt.Errorf("`start_time` is %q and should be at least %q in the future", t, duration)
		}

		parameters.Properties.SetStartTimeAsTime(t.In(loc))
	} else {
		parameters.Properties.SetStartTimeAsTime(time.Now().In(loc).Add(time.Duration(7) * time.Minute))
	}

	if v, ok := d.GetOk("expiry_time"); ok {
		parameters.Properties.ExpiryTime = pointer.To(v.(string))
	}

	// only pay attention to interval if frequency is not OneTime, and default it to 1 if not set
	if parameters.Properties.Frequency != schedule.ScheduleFrequencyOneTime {

		var interval interface{}
		interval = 1
		if v, ok := d.GetOk("interval"); ok {
			interval = v
		}
		parameters.Properties.Interval = &interval
	}

	// only pay attention to the advanced schedule fields if frequency is either Week or Month
	if parameters.Properties.Frequency == schedule.ScheduleFrequencyWeek || parameters.Properties.Frequency == schedule.ScheduleFrequencyMonth {
		parameters.Properties.AdvancedSchedule = expandArmAutomationScheduleAdvanced(d, d.Id() != "")
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationScheduleRead(d, meta)
}

func resourceAutomationScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Schedule
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schedule.ParseScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ScheduleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("frequency", string(pointer.From(props.Frequency)))

			startTime, err := props.GetStartTimeAsTime()
			if err != nil {
				return err
			}
			d.Set("start_time", startTime.Format(time.RFC3339))
			d.Set("expiry_time", pointer.From(props.ExpiryTime))

			if v := props.Interval; v != nil {
				d.Set("interval", v)
			}
			if v := props.Description; v != nil {
				d.Set("description", v)
			}
			if v := props.TimeZone; v != nil {
				d.Set("timezone", v)
			}

			if v := props.AdvancedSchedule; v != nil {
				if err := d.Set("week_days", flattenArmAutomationScheduleAdvancedWeekDays(v)); err != nil {
					return fmt.Errorf("setting `week_days`: %+v", err)
				}
				if err := d.Set("month_days", flattenArmAutomationScheduleAdvancedMonthDays(v)); err != nil {
					return fmt.Errorf("setting `month_days`: %+v", err)
				}
				if err := d.Set("monthly_occurrence", flattenArmAutomationScheduleAdvancedMonthlyOccurrences(v)); err != nil {
					return fmt.Errorf("setting `monthly_occurrence`: %+v", err)
				}
			}
		}
	}
	return nil
}

func resourceAutomationScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.Schedule
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schedule.ParseScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandArmAutomationScheduleAdvanced(d *pluginsdk.ResourceData, isUpdate bool) *schedule.AdvancedSchedule {
	expandedAdvancedSchedule := schedule.AdvancedSchedule{}

	// If frequency is set to `Month` the `week_days` array cannot be set (even empty), otherwise the API returns an error.
	// During update it can be set and it will not return an error. Workaround for the APIs behaviour
	if v, ok := d.GetOk("week_days"); ok {
		weekDays := v.(*pluginsdk.Set).List()
		expandedWeekDays := make([]string, len(weekDays))
		for i := range weekDays {
			expandedWeekDays[i] = weekDays[i].(string)
		}
		expandedAdvancedSchedule.WeekDays = &expandedWeekDays
	} else if isUpdate {
		expandedAdvancedSchedule.WeekDays = &[]string{}
	}

	// Same as above with `week_days`
	if v, ok := d.GetOk("month_days"); ok {
		monthDays := v.(*pluginsdk.Set).List()
		expandedMonthDays := make([]int64, len(monthDays))
		for i := range monthDays {
			expandedMonthDays[i] = int64(monthDays[i].(int))
		}
		expandedAdvancedSchedule.MonthDays = &expandedMonthDays
	} else if isUpdate {
		expandedAdvancedSchedule.MonthDays = &[]int64{}
	}

	monthlyOccurrences := d.Get("monthly_occurrence").([]interface{})
	expandedMonthlyOccurrences := make([]schedule.AdvancedScheduleMonthlyOccurrence, len(monthlyOccurrences))
	for i := range monthlyOccurrences {
		m := monthlyOccurrences[i].(map[string]interface{})
		occurrence := int64(m["occurrence"].(int))

		day := schedule.ScheduleDay(m["day"].(string))

		expandedMonthlyOccurrences[i] = schedule.AdvancedScheduleMonthlyOccurrence{
			Occurrence: &occurrence,
			Day:        &day,
		}
	}
	expandedAdvancedSchedule.MonthlyOccurrences = &expandedMonthlyOccurrences

	return &expandedAdvancedSchedule
}

func flattenArmAutomationScheduleAdvancedWeekDays(s *schedule.AdvancedSchedule) *pluginsdk.Set {
	flattenedWeekDays := pluginsdk.NewSet(set.HashStringIgnoreCase, []interface{}{})
	if weekDays := s.WeekDays; weekDays != nil {
		for _, v := range *weekDays {
			flattenedWeekDays.Add(v)
		}
	}
	return flattenedWeekDays
}

func flattenArmAutomationScheduleAdvancedMonthDays(s *schedule.AdvancedSchedule) *pluginsdk.Set {
	flattenedMonthDays := pluginsdk.NewSet(set.HashInt, []interface{}{})
	if monthDays := s.MonthDays; monthDays != nil {
		for _, v := range *monthDays {
			flattenedMonthDays.Add(int(v))
		}
	}
	return flattenedMonthDays
}

func flattenArmAutomationScheduleAdvancedMonthlyOccurrences(s *schedule.AdvancedSchedule) []map[string]interface{} {
	flattenedMonthlyOccurrences := make([]map[string]interface{}, 0)
	if monthlyOccurrences := s.MonthlyOccurrences; monthlyOccurrences != nil {
		for _, v := range *monthlyOccurrences {
			f := make(map[string]interface{})
			f["day"] = v.Day
			f["occurrence"] = int(*v.Occurrence)
			flattenedMonthlyOccurrences = append(flattenedMonthlyOccurrences, f)
		}
	}
	return flattenedMonthlyOccurrences
}
