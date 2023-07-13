// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryTriggerSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryTriggerScheduleCreate,
		Read:   resourceDataFactoryTriggerScheduleRead,
		Update: resourceDataFactoryTriggerScheduleUpdate,
		Delete: resourceDataFactoryTriggerScheduleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TriggerID(id)
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
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"days_of_month": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeInt,
								ValidateFunc: validation.Any(
									validation.IntBetween(1, 31),
									validation.IntBetween(-31, -1),
								),
							},
						},

						"days_of_week": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 7,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsDayOfTheWeek(false),
							},
						},

						"hours": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(0, 24),
							},
						},

						"minutes": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(0, 60),
							},
						},

						"monthly": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"weekday": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsDayOfTheWeek(false),
									},

									"week": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IntBetween(1, 5),
											validation.IntBetween(-5, -1),
										),
									},
								},
							},
						},
					},
				},
			},

			// This time can only be  represented in UTC.
			// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
			// https://github.com/Azure/azure-sdk-for-go/issues/6244
			"start_time": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
			},

			// This time can only be  represented in UTC.
			// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
			// https://github.com/Azure/azure-sdk-for-go/issues/6244
			"end_time": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
			},

			"time_zone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"frequency": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(datafactory.RecurrenceFrequencyMinute),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.RecurrenceFrequencyMinute),
					string(datafactory.RecurrenceFrequencyHour),
					string(datafactory.RecurrenceFrequencyDay),
					string(datafactory.RecurrenceFrequencyWeek),
					string(datafactory.RecurrenceFrequencyMonth),
				}, false),
			},

			"interval": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"activated": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"pipeline": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"pipeline_parameters"},
				ExactlyOneOf:  []string{"pipeline", "pipeline_name"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
						},

						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"pipeline_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"pipeline", "pipeline_name"},
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"pipeline_parameters": {
				Type:          pluginsdk.TypeMap,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"pipeline"},
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceDataFactoryTriggerScheduleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_factory_trigger_schedule", id.ID())
	}

	props := &datafactory.ScheduleTriggerTypeProperties{
		Recurrence: &datafactory.ScheduleTriggerRecurrence{
			Frequency: datafactory.RecurrenceFrequency(d.Get("frequency").(string)),
			Interval:  utils.Int32(int32(d.Get("interval").(int))),
			Schedule:  expandDataFactorySchedule(d.Get("schedule").([]interface{})),
			TimeZone:  utils.String(d.Get("time_zone").(string)),
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
		Description:                   utils.String(d.Get("description").(string)),
	}

	if pipelineName := d.Get("pipeline_name").(string); len(pipelineName) != 0 {
		scheduleProps.Pipelines = &[]datafactory.TriggerPipelineReference{
			{
				PipelineReference: &datafactory.PipelineReference{
					ReferenceName: utils.String(pipelineName),
					Type:          utils.String("PipelineReference"),
				},
				Parameters: d.Get("pipeline_parameters").(map[string]interface{}),
			},
		}
	} else {
		scheduleProps.Pipelines = expandDataFactoryPipelines(d.Get("pipeline").([]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		scheduleProps.Annotations = &annotations
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

func resourceDataFactoryTriggerScheduleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
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
			Interval:  utils.Int32(int32(d.Get("interval").(int))),
			Schedule:  expandDataFactorySchedule(d.Get("schedule").([]interface{})),
			TimeZone:  utils.String(d.Get("time_zone").(string)),
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
		Description:                   utils.String(d.Get("description").(string)),
	}

	pipelineName := d.Get("pipeline_name").(string)
	pipeline := d.Get("pipeline").([]interface{})
	if (d.HasChange("pipeline_name") && len(pipelineName) == 0) || (d.HasChange("pipeline") && len(pipeline) != 0) {
		scheduleProps.Pipelines = expandDataFactoryPipelines(pipeline)
	} else {
		scheduleProps.Pipelines = &[]datafactory.TriggerPipelineReference{
			{
				PipelineReference: &datafactory.PipelineReference{
					ReferenceName: utils.String(pipelineName),
					Type:          utils.String("PipelineReference"),
				},
				Parameters: d.Get("pipeline_parameters").(map[string]interface{}),
			},
		}
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		scheduleProps.Annotations = &annotations
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

func resourceDataFactoryTriggerScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if utils.ResponseWasNotFound(resp.Response) {
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

		annotations := flattenDataFactoryAnnotations(scheduleTriggerProps.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		d.Set("description", scheduleTriggerProps.Description)
	}

	return nil
}

func resourceDataFactoryTriggerScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("stopping %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to stop %s: %+v", id, err)
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandDataFactorySchedule(input []interface{}) *datafactory.RecurrenceSchedule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	schedule := datafactory.RecurrenceSchedule{}

	value := input[0].(map[string]interface{})
	weekDays := make([]datafactory.DaysOfWeek, 0)
	for _, v := range value["days_of_week"].([]interface{}) {
		weekDays = append(weekDays, datafactory.DaysOfWeek(v.(string)))
	}
	if len(weekDays) > 0 {
		schedule.WeekDays = &weekDays
	}

	monthlyOccurrences := make([]datafactory.RecurrenceScheduleOccurrence, 0)
	for _, v := range value["monthly"].([]interface{}) {
		value := v.(map[string]interface{})
		monthlyOccurrences = append(monthlyOccurrences, datafactory.RecurrenceScheduleOccurrence{
			Day:        datafactory.DayOfWeek(value["weekday"].(string)),
			Occurrence: utils.Int32(int32(value["week"].(int))),
		})
	}
	if len(monthlyOccurrences) > 0 {
		schedule.MonthlyOccurrences = &monthlyOccurrences
	}

	if monthdays := value["days_of_month"].([]interface{}); len(monthdays) > 0 {
		schedule.MonthDays = utils.ExpandInt32Slice(monthdays)
	}
	if minutes := value["minutes"].([]interface{}); len(minutes) > 0 {
		schedule.Minutes = utils.ExpandInt32Slice(minutes)
	}
	if hours := value["hours"].([]interface{}); len(hours) > 0 {
		schedule.Hours = utils.ExpandInt32Slice(hours)
	}

	return &schedule
}

func flattenDataFactorySchedule(schedule *datafactory.RecurrenceSchedule) []interface{} {
	if schedule == nil {
		return []interface{}{}
	}
	value := make(map[string]interface{})
	if schedule.Minutes != nil {
		value["minutes"] = utils.FlattenInt32Slice(schedule.Minutes)
	}
	if schedule.Hours != nil {
		value["hours"] = utils.FlattenInt32Slice(schedule.Hours)
	}
	if schedule.WeekDays != nil {
		weekDays := make([]interface{}, 0)
		for _, v := range *schedule.WeekDays {
			weekDays = append(weekDays, string(v))
		}
		value["days_of_week"] = weekDays
	}
	if schedule.MonthDays != nil {
		value["days_of_month"] = utils.FlattenInt32Slice(schedule.MonthDays)
	}
	if schedule.MonthlyOccurrences != nil {
		monthlyOccurrences := make([]interface{}, 0)
		for _, v := range *schedule.MonthlyOccurrences {
			occurrence := make(map[string]interface{})
			occurrence["weekday"] = string(v.Day)
			if v.Occurrence != nil {
				occurrence["week"] = *v.Occurrence
			}
			monthlyOccurrences = append(monthlyOccurrences, occurrence)
		}
		value["monthly"] = monthlyOccurrences
	}
	return []interface{}{value}
}

func expandDataFactoryPipelines(input []interface{}) *[]datafactory.TriggerPipelineReference {
	if len(input) == 0 {
		return nil
	}

	pipes := make([]datafactory.TriggerPipelineReference, 0)

	for _, item := range input {
		config := item.(map[string]interface{})
		v := datafactory.TriggerPipelineReference{
			PipelineReference: &datafactory.PipelineReference{
				ReferenceName: utils.String(config["name"].(string)),
				Type:          utils.String("PipelineReference"),
			},
			Parameters: config["parameters"].(map[string]interface{}),
		}
		pipes = append(pipes, v)
	}

	return &pipes
}

func flattenDataFactoryPipelines(pipelines *[]datafactory.TriggerPipelineReference) interface{} {
	if pipelines == nil {
		return []interface{}{}
	}

	res := make([]interface{}, 0)

	for _, item := range *pipelines {
		v := make(map[string]interface{})
		v["name"] = utils.String(*item.PipelineReference.ReferenceName)
		v["parameters"] = item.Parameters
		res = append(res, v)
	}

	return &res
}
