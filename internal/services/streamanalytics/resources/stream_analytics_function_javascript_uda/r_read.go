// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_function_javascript_uda

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsFunctionUDARead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := functions.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FunctionName)

	jobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)
	d.Set("stream_analytics_job_id", jobId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			function, ok := props.(functions.AggregateFunctionProperties)
			if !ok {
				return fmt.Errorf("converting to an Aggregate Function")
			}

			binding := function.Properties.Binding.(functions.JavaScriptFunctionBinding)

			d.Set("script", pointer.From(binding.Properties.Script))

			if err := d.Set("input", flattenStreamAnalyticsFunctionUDAInputs(function.Properties.Inputs)); err != nil {
				return fmt.Errorf("flattening `input`: %+v", err)
			}

			if err := d.Set("output", flattenStreamAnalyticsFunctionUDAOutput(function.Properties.Output)); err != nil {
				return fmt.Errorf("flattening `output`: %+v", err)
			}
		}
	}
	return nil
}

func flattenStreamAnalyticsFunctionUDAInputs(input *[]functions.FunctionInput) []any {
	if input == nil {
		return []any{}
	}

	outputs := make([]any, 0)

	for _, v := range *input {
		variableType := pointer.From(v.DataType)

		isConfigurationParameter := pointer.From(v.IsConfigurationParameter)

		outputs = append(outputs, map[string]any{
			"type":                    variableType,
			"configuration_parameter": isConfigurationParameter,
		})
	}

	return outputs
}

func flattenStreamAnalyticsFunctionUDAOutput(input *functions.FunctionOutput) []any {
	if input == nil {
		return []any{}
	}

	variableType := pointer.From(input.DataType)

	return []any{
		map[string]any{
			"type": variableType,
		},
	}
}
