// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_schedule

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerScheduleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	scheduleTriggerProps, ok := resp.Properties.AsScheduleTrigger()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q Received: %q", *id, datafactory.TypeBasicTriggerTypeScheduleTrigger, *resp.Type)
	}

	if scheduleTriggerProps != nil {
		d.Set("activated", scheduleTriggerProps.RuntimeState == datafactory.TriggerRuntimeStateStarted)

		if recurrence := scheduleTriggerProps.Recurrence; recurrence != nil {
			if v := recurrence.StartTime; v != nil {
				d.Set("start_time", v.Format(time.RFC3339))
			}
			if v := recurrence.EndTime; v != nil {
				d.Set("end_time", v.Format(time.RFC3339))
			}
			d.Set("frequency", recurrence.Frequency)
			d.Set("interval", recurrence.Interval)
			d.Set("time_zone", recurrence.TimeZone)

			if schedule := recurrence.Schedule; schedule != nil {
				d.Set("schedule", flattenDataFactorySchedule(schedule))
			}
		}

		if pipelines := scheduleTriggerProps.Pipelines; pipelines != nil {
			if len(*pipelines) > 0 {
				pipeline := *pipelines
				if reference := pipeline[0].PipelineReference; reference != nil {
					d.Set("pipeline_name", reference.ReferenceName)
				}
				d.Set("pipeline_parameters", pipeline[0].Parameters)
				d.Set("pipeline", flattenDataFactoryPipelines(pipelines))
			}
		}

		if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(scheduleTriggerProps.Annotations)); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		d.Set("description", scheduleTriggerProps.Description)
	}

	return nil
}

func flattenDataFactorySchedule(schedule *datafactory.RecurrenceSchedule) []any {
	if schedule == nil {
		return []any{}
	}
	value := make(map[string]any)
	if schedule.Minutes != nil {
		value["minutes"] = pluginsdk.FlattenSlice(schedule.Minutes)
	}
	if schedule.Hours != nil {
		value["hours"] = pluginsdk.FlattenSlice(schedule.Hours)
	}
	if schedule.WeekDays != nil {
		weekDays := make([]any, 0)
		for _, v := range *schedule.WeekDays {
			weekDays = append(weekDays, string(v))
		}
		value["days_of_week"] = weekDays
	}
	if schedule.MonthDays != nil {
		value["days_of_month"] = pluginsdk.FlattenSlice(schedule.MonthDays)
	}
	if schedule.MonthlyOccurrences != nil {
		monthlyOccurrences := make([]any, 0)
		for _, v := range *schedule.MonthlyOccurrences {
			occurrence := make(map[string]any)
			occurrence["weekday"] = string(v.Day)
			if v.Occurrence != nil {
				occurrence["week"] = *v.Occurrence
			}
			monthlyOccurrences = append(monthlyOccurrences, occurrence)
		}
		value["monthly"] = monthlyOccurrences
	}
	return []any{value}
}

func flattenDataFactoryPipelines(pipelines *[]datafactory.TriggerPipelineReference) any {
	if pipelines == nil {
		return []any{}
	}

	res := make([]any, 0)

	for _, item := range *pipelines {
		v := make(map[string]any)
		v["name"] = pointer.To(*item.PipelineReference.ReferenceName)
		v["parameters"] = item.Parameters
		res = append(res, v)
	}

	return &res
}
