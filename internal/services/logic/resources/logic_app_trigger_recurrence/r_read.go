// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_recurrence

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerRecurrenceRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	t, app, _, err := helpers.RetrieveLogicAppTrigger(d, meta, *id)
	if err != nil {
		return err
	}

	if t == nil {
		log.Printf("[DEBUG] Logic App %q (Resource Group %q) does not contain %s - removing from state", id.WorkflowName, id.ResourceGroupName, id.ID())
		d.SetId("")
		return nil
	}

	trigger := *t

	d.Set("name", id.TriggerName)
	d.Set("logic_app_id", app.Id)

	v := trigger["recurrence"]
	if v == nil {
		return fmt.Errorf("`recurrence` was nil for HTTP Trigger %s", id)
	}

	recurrence, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("parsing `recurrence` for HTTP Trigger %s", id)
	}

	if frequency := recurrence["frequency"]; frequency != nil {
		d.Set("frequency", frequency.(string))
	}

	if interval := recurrence["interval"]; interval != nil {
		d.Set("interval", int(interval.(float64)))
	}

	if startTime := recurrence["startTime"]; startTime != nil {
		d.Set("start_time", startTime.(string))
	}

	if timeZone := recurrence["timeZone"]; timeZone != nil {
		d.Set("time_zone", timeZone.(string))
	}

	if schedule := recurrence["schedule"]; schedule != nil {
		d.Set("schedule", flattenLogicAppTriggerRecurrenceSchedule(schedule.(map[string]any)))
	}

	return nil
}

func flattenLogicAppTriggerRecurrenceSchedule(input map[string]any) []any {
	attrs := make(map[string]any)

	if hours := input["hours"]; hours != nil {
		attrs["at_these_hours"] = hours
	}
	if minutes := input["minutes"]; minutes != nil {
		attrs["at_these_minutes"] = minutes
	}
	if days := input["weekDays"]; days != nil {
		attrs["on_these_days"] = days
	}

	return []any{attrs}
}
