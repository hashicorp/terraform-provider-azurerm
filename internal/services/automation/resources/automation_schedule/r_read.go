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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationScheduleRead(d *pluginsdk.ResourceData, meta any) error {
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
			d.Set("frequency", pointer.FromEnum(props.Frequency))

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
