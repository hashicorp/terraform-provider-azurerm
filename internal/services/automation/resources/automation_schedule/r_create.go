// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_schedule

import (
	"fmt"
	"time"

	// import time/tzdata to embed timezone information in the program
	// add this to resolve https://github.com/hashicorp/terraform-provider-azurerm/issues/20690
	_ "time/tzdata"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationScheduleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Schedule
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schedule.NewScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
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
		var interval any
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
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationScheduleRead(d, meta)
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

	monthlyOccurrences := d.Get("monthly_occurrence").([]any)
	expandedMonthlyOccurrences := make([]schedule.AdvancedScheduleMonthlyOccurrence, len(monthlyOccurrences))
	for i := range monthlyOccurrences {
		m := monthlyOccurrences[i].(map[string]any)

		expandedMonthlyOccurrences[i] = schedule.AdvancedScheduleMonthlyOccurrence{
			Occurrence: pointer.To(int64(m["occurrence"].(int))),
			Day:        pointer.ToEnum[schedule.ScheduleDay](m["day"].(string)),
		}
	}
	expandedAdvancedSchedule.MonthlyOccurrences = &expandedMonthlyOccurrences

	return &expandedAdvancedSchedule
}
