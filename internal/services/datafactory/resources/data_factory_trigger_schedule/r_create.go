// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_schedule

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerScheduleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_schedule", id.ID())
		}
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

	if pipelineName := d.Get("pipeline_name").(string); len(pipelineName) != 0 {
		scheduleProps.Pipelines = &[]datafactory.TriggerPipelineReference{
			{
				PipelineReference: &datafactory.PipelineReference{
					ReferenceName: pointer.To(pipelineName),
					Type:          pointer.To("PipelineReference"),
				},
				Parameters: d.Get("pipeline_parameters").(map[string]any),
			},
		}
	} else {
		scheduleProps.Pipelines = expandDataFactoryPipelines(d.Get("pipeline").([]any))
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

	d.SetId(id.ID())

	return resourceDataFactoryTriggerScheduleRead(d, meta)
}

func expandDataFactorySchedule(input []any) *datafactory.RecurrenceSchedule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	schedule := datafactory.RecurrenceSchedule{}

	value := input[0].(map[string]any)
	weekDays := make([]datafactory.DaysOfWeek, 0)
	for _, v := range value["days_of_week"].([]any) {
		weekDays = append(weekDays, datafactory.DaysOfWeek(v.(string)))
	}
	if len(weekDays) > 0 {
		schedule.WeekDays = &weekDays
	}

	monthlyOccurrences := make([]datafactory.RecurrenceScheduleOccurrence, 0)
	for _, v := range value["monthly"].([]any) {
		value := v.(map[string]any)
		monthlyOccurrences = append(monthlyOccurrences, datafactory.RecurrenceScheduleOccurrence{
			Day:        datafactory.DayOfWeek(value["weekday"].(string)),
			Occurrence: pointer.To(int32(value["week"].(int))),
		})
	}
	if len(monthlyOccurrences) > 0 {
		schedule.MonthlyOccurrences = &monthlyOccurrences
	}

	if monthdays := value["days_of_month"].([]any); len(monthdays) > 0 {
		schedule.MonthDays = pluginsdk.ExpandInt32Slice(monthdays)
	}
	if minutes := value["minutes"].([]any); len(minutes) > 0 {
		schedule.Minutes = pluginsdk.ExpandInt32Slice(minutes)
	}
	if hours := value["hours"].([]any); len(hours) > 0 {
		schedule.Hours = pluginsdk.ExpandInt32Slice(hours)
	}

	return &schedule
}

func expandDataFactoryPipelines(input []any) *[]datafactory.TriggerPipelineReference {
	if len(input) == 0 {
		return nil
	}

	pipes := make([]datafactory.TriggerPipelineReference, 0)

	for _, item := range input {
		config := item.(map[string]any)
		v := datafactory.TriggerPipelineReference{
			PipelineReference: &datafactory.PipelineReference{
				ReferenceName: pointer.To(config["name"].(string)),
				Type:          pointer.To("PipelineReference"),
			},
			Parameters: config["parameters"].(map[string]any),
		}
		pipes = append(pipes, v)
	}

	return &pipes
}
