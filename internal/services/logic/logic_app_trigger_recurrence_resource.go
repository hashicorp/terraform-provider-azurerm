// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceLogicAppTriggerRecurrence() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppTriggerRecurrenceCreateUpdate,
		Read:   resourceLogicAppTriggerRecurrenceRead,
		Update: resourceLogicAppTriggerRecurrenceCreateUpdate,
		Delete: resourceLogicAppTriggerRecurrenceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workflowtriggers.ParseTriggerID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"logic_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workflows.ValidateWorkflowID,
			},

			"frequency": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Month",
					"Week",
					"Day",
					"Hour",
					"Minute",
					"Hour",
					"Second",
				}, false),
			},

			"interval": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"start_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"at_these_hours": {
							Type:         pluginsdk.TypeSet,
							Optional:     true,
							AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(0, 23),
							},
						},
						"at_these_minutes": {
							Type:         pluginsdk.TypeSet,
							Optional:     true,
							AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(0, 59),
							},
						},
						"on_these_days": {
							Type:         pluginsdk.TypeSet,
							Optional:     true,
							AtLeastOneOf: []string{"schedule.0.at_these_hours", "schedule.0.at_these_minutes", "schedule.0.on_these_days"},
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

			"time_zone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.TriggerRecurrenceTimeZone(),
			},
		},
	}
}

func resourceLogicAppTriggerRecurrenceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	trigger := map[string]interface{}{
		"recurrence": map[string]interface{}{
			"frequency": d.Get("frequency").(string),
			"interval":  d.Get("interval").(int),
		},
		"type": "Recurrence",
	}

	if v, ok := d.GetOk("start_time"); ok {
		trigger["recurrence"].(map[string]interface{})["startTime"] = v.(string)
	}

	if v, ok := d.GetOk("time_zone"); ok {
		trigger["recurrence"].(map[string]interface{})["timeZone"] = v.(string)
	}

	if v, ok := d.GetOk("schedule"); ok {
		trigger["recurrence"].(map[string]interface{})["schedule"] = expandLogicAppTriggerRecurrenceSchedule(v.([]interface{}))
	}

	workflowId, err := workflows.ParseWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	id := workflowtriggers.NewTriggerID(workflowId.SubscriptionId, workflowId.ResourceGroupName, workflowId.WorkflowName, d.Get("name").(string))

	if err := resourceLogicAppTriggerUpdate(d, meta, *workflowId, id, trigger, "azurerm_logic_app_trigger_recurrence"); err != nil {
		return err
	}

	return resourceLogicAppTriggerRecurrenceRead(d, meta)
}

func resourceLogicAppTriggerRecurrenceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroupName, id.WorkflowName)

	t, app, err := retrieveLogicAppTrigger(d, meta, workflowId, id.TriggerName)
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

	recurrence, ok := v.(map[string]interface{})
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
		d.Set("schedule", flattenLogicAppTriggerRecurrenceSchedule(schedule.(map[string]interface{})))
	}

	return nil
}

func resourceLogicAppTriggerRecurrenceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := workflowtriggers.ParseTriggerID(d.Id())
	if err != nil {
		return err
	}

	workflowId := workflows.NewWorkflowID(id.SubscriptionId, id.ResourceGroupName, id.WorkflowName)

	err = resourceLogicAppTriggerRemove(d, meta, workflowId, id.TriggerName)
	if err != nil {
		return fmt.Errorf("removing Trigger %s: %+v", id, err)
	}

	return nil
}

func expandLogicAppTriggerRecurrenceSchedule(input []interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	if len(input) == 0 || input[0] == nil {
		return output
	}

	attrs := input[0].(map[string]interface{})
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

func flattenLogicAppTriggerRecurrenceSchedule(input map[string]interface{}) []interface{} {
	attrs := make(map[string]interface{})

	if hours := input["hours"]; hours != nil {
		attrs["at_these_hours"] = hours
	}
	if minutes := input["minutes"]; minutes != nil {
		attrs["at_these_minutes"] = minutes
	}
	if days := input["weekDays"]; days != nil {
		attrs["on_these_days"] = days
	}

	return []interface{}{attrs}
}
