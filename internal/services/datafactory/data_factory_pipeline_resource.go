// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryPipeline() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryPipelineCreateUpdate,
		Read:   resourceDataFactoryPipelineRead,
		Update: resourceDataFactoryPipelineCreateUpdate,
		Delete: resourceDataFactoryPipelineDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := pipelines.ParsePipelineID(id)
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

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"variables": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"activities_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"concurrency": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"moniter_metrics_after_duration": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDataFactoryPipelineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := pipelines.NewPipelineID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, pipelines.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory_pipeline", id.ID())
		}
	}

	payload := pipelines.PipelineResource{
		Properties: pipelines.Pipeline{
			Description: pointer.To(d.Get("description").(string)),
			Parameters:  expandDataFactoryPipelineParameters(d.Get("parameters").(map[string]interface{})),
			Variables:   expandDataFactoryPipelineVariables(d.Get("variables").(map[string]interface{})),
		},
	}

	if v, ok := d.GetOk("activities_json"); ok && v.(string) != "" {
		activitiesJson, err := pipelines.UnmarshalActivityImplementation([]byte(fmt.Sprintf(`{ "activities": %s }`, v.(string))))
		if err != nil {
			return fmt.Errorf("unmarshaling `activities_json`: %+v", err)
		}
		rawActivities, ok := activitiesJson.(pipelines.RawActivityImpl)
		if !ok {
			return fmt.Errorf("expected `activities_json` to be of type `RawActivityImpl`")
		}

		activities := make([]pipelines.Activity, 0)
		acts, ok := rawActivities.Values["activities"]
		if !ok {
			return fmt.Errorf("`activities` was not found in the unmarshaled `activities_json`")
		}

		for _, activity := range acts.([]interface{}) {
			act, err := json.Marshal(activity)
			if err != nil {
				return fmt.Errorf("marshaling activity %+v: %+v", activity, err)
			}
			a, err := pipelines.UnmarshalActivityImplementation(act)
			if err != nil {
				return fmt.Errorf("unmarshaling activity %+v: %+v", act, err)
			}
			activities = append(activities, a)
		}

		payload.Properties.Activities = pointer.To(activities)
	}

	annotations := make([]interface{}, 0)
	if v, ok := d.GetOk("annotations"); ok {
		annotations = v.([]interface{})
	}
	payload.Properties.Annotations = &annotations

	if v, ok := d.GetOk("concurrency"); ok {
		payload.Properties.Concurrency = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("moniter_metrics_after_duration"); ok {
		payload.Properties.Policy = &pipelines.PipelinePolicy{
			ElapsedTimeMetric: &pipelines.PipelineElapsedTimeMetricPolicy{
				Duration: pointer.To(v),
			},
		}
	}

	if v, ok := d.GetOk("folder"); ok {
		payload.Properties.Folder = &pipelines.PipelineFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload, pipelines.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryPipelineRead(d, meta)
}

func resourceDataFactoryPipelineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pipelines.ParsePipelineID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, pipelines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PipelineName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("description", pointer.From(props.Description))

		parameters := flattenDataFactoryPipelineParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		annotations := flattenDataFactoryAnnotations(props.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		d.Set("concurrency", pointer.From(props.Concurrency))

		elapsedTimeMetricDuration := ""
		if props.Policy != nil && props.Policy.ElapsedTimeMetric != nil && props.Policy.ElapsedTimeMetric.Duration != nil {
			if v, ok := (*props.Policy.ElapsedTimeMetric.Duration).(string); ok {
				elapsedTimeMetricDuration = v
			}
		}
		d.Set("moniter_metrics_after_duration", elapsedTimeMetricDuration)

		if folder := props.Folder; folder != nil {
			d.Set("folder", pointer.From(folder.Name))
		}

		variables := flattenDataFactoryPipelineVariables(props.Variables)
		if err := d.Set("variables", variables); err != nil {
			return fmt.Errorf("setting `variables`: %+v", err)
		}

		activitiesJson := ""
		if activities := props.Activities; activities != nil {
			acts, err := json.Marshal(activities)
			if err != nil {
				return fmt.Errorf("marshaling `activities_json`: %+v", err)
			}

			activitiesJson = string(acts)
		}
		d.Set("activities_json", activitiesJson)
	}

	return nil
}

func resourceDataFactoryPipelineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pipelines.ParsePipelineID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandDataFactoryPipelineParameters(input map[string]interface{}) *map[string]pipelines.ParameterSpecification {
	output := make(map[string]pipelines.ParameterSpecification)

	for k, v := range input {
		output[k] = pipelines.ParameterSpecification{
			Type:         pipelines.ParameterTypeString,
			DefaultValue: pointer.To(v),
		}
	}

	return &output
}

func flattenDataFactoryPipelineParameters(input *map[string]pipelines.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}
	for k, v := range *input {
		// we only support string parameters at this time
		if v.Type != pipelines.ParameterTypeString {
			log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			continue
		}

		output[k] = pointer.From(v.DefaultValue)
	}

	return output
}

func expandDataFactoryPipelineVariables(input map[string]interface{}) *map[string]pipelines.VariableSpecification {
	output := make(map[string]pipelines.VariableSpecification)

	for k, v := range input {
		output[k] = pipelines.VariableSpecification{
			Type:         pipelines.VariableTypeString,
			DefaultValue: pointer.To(v),
		}
	}

	return &output
}

func flattenDataFactoryPipelineVariables(input *map[string]pipelines.VariableSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	for k, v := range *input {
		if v.Type != pipelines.VariableTypeString {
			log.Printf("[DEBUG] Skipping variable %q since it's not a string", k)
			continue
		}

		output[k] = pointer.From(v.DefaultValue)
	}

	return output
}
