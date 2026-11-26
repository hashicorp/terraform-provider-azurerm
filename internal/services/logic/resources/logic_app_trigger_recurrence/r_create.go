// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_trigger_recurrence

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceLogicAppTriggerRecurrenceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	trigger := map[string]any{
		"recurrence": map[string]any{
			"frequency": d.Get("frequency").(string),
			"interval":  d.Get("interval").(int),
		},
		"type": "Recurrence",
	}

	if v, ok := d.GetOk("start_time"); ok {
		trigger["recurrence"].(map[string]any)["startTime"] = v.(string)
	}

	if v, ok := d.GetOk("time_zone"); ok {
		trigger["recurrence"].(map[string]any)["timeZone"] = v.(string)
	}

	if v, ok := d.GetOk("schedule"); ok {
		trigger["recurrence"].(map[string]any)["schedule"] = expandLogicAppTriggerRecurrenceSchedule(v.([]any))
	}

	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := workflowtriggers.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	if err := helpers.ResourceLogicAppTriggerUpdate(d, meta, *workflowId, id, trigger, "azurerm_logic_app_trigger_recurrence"); err != nil {
		return err
	}

	return resourceLogicAppTriggerRecurrenceRead(d, meta)
}

func expandLogicAppTriggerRecurrenceSchedule(input []any) map[string]any {
	output := make(map[string]any)
	if len(input) == 0 || input[0] == nil {
		return output
	}

	attrs := input[0].(map[string]any)
	if hoursRaw, ok := attrs["at_these_hours"]; ok {
		hoursSet := hoursRaw.(*pluginsdk.Set).List()
		hours := make([]int, 0)
		for _, hour := range hoursSet {
			hours = append(hours, hour.(int))
		}
		if len(hours) > 0 {
			output["hours"] = &hours
		}
	}
	if minutesRaw, ok := attrs["at_these_minutes"]; ok {
		minutesSet := minutesRaw.(*pluginsdk.Set).List()
		minutes := make([]int, 0)
		for _, minute := range minutesSet {
			minutes = append(minutes, minute.(int))
		}
		if len(minutes) > 0 {
			output["minutes"] = &minutes
		}
	}
	if daysRaw, ok := attrs["on_these_days"]; ok {
		daysSet := daysRaw.(*pluginsdk.Set).List()
		days := make([]string, 0)
		for _, day := range daysSet {
			days = append(days, day.(string))
		}
		if len(days) > 0 {
			output["weekDays"] = &days
		}
	}

	return output
}
