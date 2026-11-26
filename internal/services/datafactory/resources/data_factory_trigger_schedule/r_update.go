// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_schedule

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerScheduleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, ""); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// activated triggers cannot be updated - we activate the trigger again after updating
	future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("stopping %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to stop %s: %+v", id, err)
	}

	props := &datafactory.ScheduleTriggerTypeProperties{
		Recurrence: &datafactory.ScheduleTriggerRecurrence{
			Frequency: datafactory.RecurrenceFrequency(d.Get("frequency").(string)),
			Interval:  pointer.To(int32(d.Get("interval").(int))),
			Schedule:  expandDataFactorySchedule(d.Get("schedule").([]any)),
			TimeZone:  pointer.To(d.Get("time_zone").(string)),
		},
	}

	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.Recurrence.StartTime = &date.Time{Time: t}
	} else {
		t, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
		props.Recurrence.StartTime = &date.Time{Time: t}
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.Recurrence.EndTime = &date.Time{Time: t}
	}

	scheduleProps := &datafactory.ScheduleTrigger{
		ScheduleTriggerTypeProperties: props,
		Description:                   pointer.To(d.Get("description").(string)),
	}

	pipelineName := d.Get("pipeline_name").(string)
	pipeline := d.Get("pipeline").([]any)
	if (d.HasChange("pipeline_name") && len(pipelineName) == 0) || (d.HasChange("pipeline") && len(pipeline) != 0) {
		scheduleProps.Pipelines = expandDataFactoryPipelines(pipeline)
	} else {
		scheduleProps.Pipelines = &[]datafactory.TriggerPipelineReference{
			{
				PipelineReference: &datafactory.PipelineReference{
					ReferenceName: pointer.To(pipelineName),
					Type:          pointer.To("PipelineReference"),
				},
				Parameters: d.Get("pipeline_parameters").(map[string]any),
			},
		}
	}

	if v, ok := d.GetOk("annotations"); ok {
		scheduleProps.Annotations = pointer.To(v.([]any))
	}

	trigger := datafactory.TriggerResource{
		Properties: scheduleProps,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, trigger, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.Get("activated").(bool) {
		future, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on start %s: %+v", id, err)
		}
	}

	return resourceDataFactoryTriggerScheduleRead(d, meta)
}
