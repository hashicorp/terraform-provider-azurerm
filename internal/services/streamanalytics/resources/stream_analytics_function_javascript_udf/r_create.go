// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_function_javascript_udf

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsFunctionUDFCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := functions.NewFunctionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_function_javascript_udf", id.ID())
			}
		}
	}

	inputs := expandStreamAnalyticsFunctionInputs(d.Get("input").([]any))
	output := expandStreamAnalyticsFunctionOutput(d.Get("output").([]any))

	function := functions.Function{
		Properties: &functions.ScalarFunctionProperties{
			Properties: &functions.FunctionConfiguration{
				Binding: &functions.JavaScriptFunctionBinding{
					Properties: &functions.JavaScriptFunctionBindingProperties{
						Script: pointer.To(d.Get("script").(string)),
					},
				},
				Inputs: inputs,
				Output: output,
			},
		},
	}

	var createOpts functions.CreateOrReplaceOperationOptions
	var updateOpts functions.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, function, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, function, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsFunctionUDFRead(d, meta)
}

func expandStreamAnalyticsFunctionInputs(input []any) *[]functions.FunctionInput {
	outputs := make([]functions.FunctionInput, 0)

	for _, raw := range input {
		v := raw.(map[string]any)
		variableType := v["type"].(string)
		outputs = append(outputs, functions.FunctionInput{
			DataType:                 pointer.To(variableType),
			IsConfigurationParameter: pointer.To(v["configuration_parameter"].(bool)),
		})
	}

	return &outputs
}

func expandStreamAnalyticsFunctionOutput(input []any) *functions.FunctionOutput {
	output := input[0].(map[string]any)

	dataType := output["type"].(string)
	return &functions.FunctionOutput{
		DataType: pointer.To(dataType),
	}
}
