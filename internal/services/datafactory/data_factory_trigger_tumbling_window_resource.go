// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

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
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryTriggerTumblingWindow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryTriggerTumblingWindowCreateUpdate,
		Read:   resourceDataFactoryTriggerTumblingWindowRead,
		Update: resourceDataFactoryTriggerTumblingWindowCreateUpdate,
		Delete: resourceDataFactoryTriggerTumblingWindowDelete,

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

			"frequency": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.TumblingWindowFrequencyHour),
					string(datafactory.TumblingWindowFrequencyMinute),
					string(datafactory.TumblingWindowFrequencyMonth),
				}, false),
			},

			"interval": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"pipeline": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
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

			"start_time": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"activated": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
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

			"delay": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.TriggerTimespan,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"end_time": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"max_concurrency": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validation.IntBetween(1, 50),
			},

			"retry": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"count": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},

						"interval": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"trigger_dependency": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"offset": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.TriggerTimespan,
						},

						"size": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.TriggerTimespan,
						},

						"trigger_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDataFactoryTriggerTumblingWindowCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_tumbling_window", id.ID())
		}
	} else {
		future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to stop %s: %+v", id, err)
		}
	}

	startTime, err := time.Parse(time.RFC3339, d.Get("start_time").(string))
	if err != nil {
		return err
	}

	props := &datafactory.TumblingWindowTrigger{
		TumblingWindowTriggerTypeProperties: &datafactory.TumblingWindowTriggerTypeProperties{
			Frequency:      datafactory.TumblingWindowFrequency(d.Get("frequency").(string)),
			Interval:       utils.Int32(int32(d.Get("interval").(int))),
			MaxConcurrency: utils.Int32(int32(d.Get("max_concurrency").(int))),
			RetryPolicy:    expandDataFactoryTriggerTumblingWindowRetryPolicy(d.Get("retry").([]interface{})),
			DependsOn:      expandDataFactoryTriggerDependency(d.Get("trigger_dependency").(*pluginsdk.Set).List()),
			StartTime:      &date.Time{Time: startTime},
		},
		Description: utils.String(d.Get("description").(string)),
		Pipeline:    expandDataFactoryTriggerSinglePipeline(d.Get("pipeline").([]interface{})),
		Type:        datafactory.TypeBasicTriggerTypeTumblingWindowTrigger,
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return err
		}
		props.TumblingWindowTriggerTypeProperties.EndTime = &date.Time{Time: t}
	}

	if v, ok := d.GetOk("delay"); ok {
		props.Delay = v.(string)
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		props.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		props.AdditionalProperties = v.(map[string]interface{})
	}

	trigger := datafactory.TriggerResource{
		Properties: props,
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

	return resourceDataFactoryTriggerTumblingWindowRead(d, meta)
}

func resourceDataFactoryTriggerTumblingWindowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	trigger, ok := resp.Properties.AsTumblingWindowTrigger()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q", id, datafactory.TypeBasicTriggerTypeTumblingWindowTrigger)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	d.Set("activated", trigger.RuntimeState == datafactory.TriggerRuntimeStateStarted)
	d.Set("additional_properties", trigger.AdditionalProperties)
	d.Set("description", trigger.Description)

	if err := d.Set("annotations", flattenDataFactoryAnnotations(trigger.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}
	if err := d.Set("pipeline", flattenDataFactoryTriggerSinglePipeline(trigger.Pipeline)); err != nil {
		return fmt.Errorf("setting `pipeline`: %+v", err)
	}

	if props := trigger.TumblingWindowTriggerTypeProperties; props != nil {
		d.Set("frequency", props.Frequency)

		interval := 0
		if props.Interval != nil {
			interval = int(*props.Interval)
		}
		d.Set("interval", interval)

		maxConcurrency := 0
		if props.MaxConcurrency != nil {
			maxConcurrency = int(*props.MaxConcurrency)
		}
		d.Set("max_concurrency", maxConcurrency)

		startTime := ""
		if v := props.StartTime; v != nil {
			startTime = v.Format(time.RFC3339)
		}
		d.Set("start_time", startTime)

		endTime := ""
		if v := props.EndTime; v != nil {
			endTime = v.Format(time.RFC3339)
		}
		d.Set("end_time", endTime)

		delay := ""
		if v, ok := props.Delay.(string); ok {
			delay = v
		}
		d.Set("delay", delay)

		if err := d.Set("retry", flattenDataFactoryTriggerRetryPolicy(props.RetryPolicy)); err != nil {
			return fmt.Errorf("setting `retry`: %+v", err)
		}

		if err := d.Set("trigger_dependency", flattenDataFactoryTriggerDependency(props.DependsOn)); err != nil {
			return fmt.Errorf("setting `trigger_dependency`: %+v", err)
		}
	}

	return nil
}

func resourceDataFactoryTriggerTumblingWindowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDataFactoryTriggerTumblingWindowRetryPolicy(input []interface{}) *datafactory.RetryPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &datafactory.RetryPolicy{
		Count:             utils.Int32(int32(raw["count"].(int))),
		IntervalInSeconds: utils.Int32(int32(raw["interval"].(int))),
	}
}

func expandDataFactoryTriggerSinglePipeline(input []interface{}) *datafactory.TriggerPipelineReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &datafactory.TriggerPipelineReference{
		PipelineReference: &datafactory.PipelineReference{
			ReferenceName: utils.String(raw["name"].(string)),
			Type:          utils.String("PipelineReference"),
		},
		Parameters: raw["parameters"].(map[string]interface{}),
	}
}

func expandDataFactoryTriggerDependency(input []interface{}) *[]datafactory.BasicDependencyReference {
	if len(input) == 0 {
		return nil
	}

	var result []datafactory.BasicDependencyReference
	for _, item := range input {
		raw := item.(map[string]interface{})

		var trigger datafactory.BasicDependencyReference

		var offset, size *string
		if v := raw["offset"].(string); v != "" {
			offset = utils.String(v)
		}
		if v := raw["size"].(string); v != "" {
			size = utils.String(v)
		}

		if v := raw["trigger_name"].(string); v != "" {
			trigger = &datafactory.TumblingWindowTriggerDependencyReference{
				Offset: offset,
				Size:   size,
				ReferenceTrigger: &datafactory.TriggerReference{
					ReferenceName: utils.String(v),
					Type:          utils.String("TriggerReference"),
				},
			}
		} else {
			trigger = &datafactory.SelfDependencyTumblingWindowTriggerReference{
				Offset: offset,
				Size:   size,
			}
		}

		result = append(result, trigger)
	}
	return &result
}

func flattenDataFactoryTriggerRetryPolicy(input *datafactory.RetryPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	count := 0
	// a little wield: after tested, it's of type float64
	if v, ok := input.Count.(float64); ok {
		count = int(v)
	}

	interval := 0
	if input.IntervalInSeconds != nil {
		interval = int(*input.IntervalInSeconds)
	}

	return []interface{}{
		map[string]interface{}{
			"count":    count,
			"interval": interval,
		},
	}
}

func flattenDataFactoryTriggerSinglePipeline(input *datafactory.TriggerPipelineReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.PipelineReference != nil && input.PipelineReference.ReferenceName != nil {
		name = *input.PipelineReference.ReferenceName
	}

	return []interface{}{
		map[string]interface{}{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}

func flattenDataFactoryTriggerDependency(input *[]datafactory.BasicDependencyReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		var offset, size, triggerName string

		switch item := item.(type) {
		case datafactory.TumblingWindowTriggerDependencyReference:
			if item.Size != nil {
				size = *item.Size
			}
			if item.Offset != nil {
				offset = *item.Offset
			}
			if item.ReferenceTrigger != nil && item.ReferenceTrigger.ReferenceName != nil {
				triggerName = *item.ReferenceTrigger.ReferenceName
			}
		case datafactory.SelfDependencyTumblingWindowTriggerReference:
			if item.Size != nil {
				size = *item.Size
			}
			if item.Offset != nil {
				offset = *item.Offset
			}
		}

		result = append(result, map[string]interface{}{
			"trigger_name": triggerName,
			"offset":       offset,
			"size":         size,
		})
	}
	return result
}
