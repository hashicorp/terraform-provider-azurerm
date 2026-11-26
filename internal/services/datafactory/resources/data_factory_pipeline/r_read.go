// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_pipeline

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryPipelineRead(d *pluginsdk.ResourceData, meta any) error {
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

		if err := d.Set("parameters", flattenDataFactoryPipelineParameters(props.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(props.Annotations)); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		d.Set("concurrency", pointer.From(props.Concurrency))

		elapsedTimeMetricDuration := ""
		if props.Policy != nil && props.Policy.ElapsedTimeMetric != nil && props.Policy.ElapsedTimeMetric.Duration != nil {
			if v, ok := (*props.Policy.ElapsedTimeMetric.Duration).(string); ok {
				elapsedTimeMetricDuration = v
			}
		}

		d.Set("monitor_metrics_after_duration", elapsedTimeMetricDuration)

		if folder := props.Folder; folder != nil {
			d.Set("folder", pointer.From(folder.Name))
		}

		if err := d.Set("variables", flattenDataFactoryPipelineVariables(props.Variables)); err != nil {
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

func flattenDataFactoryPipelineParameters(input *map[string]pipelines.ParameterSpecification) map[string]any {
	output := make(map[string]any)

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

func flattenDataFactoryPipelineVariables(input *map[string]pipelines.VariableSpecification) map[string]any {
	output := make(map[string]any)

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
