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
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationScheduleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Schedule
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schedule.NewScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %v", id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("reading existing properties of %s", id)
	}

	parameters := schedule.ScheduleCreateOrUpdateParameters{
		Name: id.ScheduleName,
		Properties: schedule.ScheduleCreateOrUpdateProperties{
			Frequency:   pointer.From(existing.Model.Properties.Frequency),
			Description: existing.Model.Properties.Description,
			TimeZone:    existing.Model.Properties.TimeZone,
		},
	}
	if d.HasChange("frequency") {
		parameters.Properties.Frequency = schedule.ScheduleFrequency(d.Get("frequency").(string))
	}

	if d.HasChange("timezone") {
		parameters.Properties.TimeZone = pointer.To(d.Get("timezone").(string))
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	parameters.Properties.StartTime = pointer.From(existing.Model.Properties.StartTime)
	if d.HasChange("start_time") {
		// start time can default to now + 7 (5 could be invalid by the time the API is called)
		loc, err := time.LoadLocation(pointer.From(parameters.Properties.TimeZone))
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
	}

	parameters.Properties.ExpiryTime = existing.Model.Properties.ExpiryTime
	if d.HasChange("expiry_time") {
		parameters.Properties.ExpiryTime = pointer.To(d.Get("expiry_time").(string))
	}

	// only pay attention to interval if frequency is not OneTime, and default it to 1 if not set
	if parameters.Properties.Frequency != schedule.ScheduleFrequencyOneTime {
		parameters.Properties.Interval = existing.Model.Properties.Interval
		if d.HasChange("interval") {
			parameters.Properties.Interval = pointer.To(d.Get("interval"))
		}
	}

	if d.HasChanges("week_days", "month_days", "monthly_occurrence") {
		// only pay attention to the advanced schedule fields if frequency is either Week or Month
		if parameters.Properties.Frequency == schedule.ScheduleFrequencyWeek || parameters.Properties.Frequency == schedule.ScheduleFrequencyMonth {
			parameters.Properties.AdvancedSchedule = expandArmAutomationScheduleAdvanced(d, d.Id() != "")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceAutomationScheduleRead(d, meta)
}
