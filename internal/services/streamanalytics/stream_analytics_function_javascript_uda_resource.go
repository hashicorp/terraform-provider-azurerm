// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsFunctionUDA() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsFunctionUDACreate,
		Read:   resourceStreamAnalyticsFunctionUDARead,
		Update: resourceStreamAnalyticsFunctionUDAUpdate,
		Delete: resourceStreamAnalyticsFunctionUDADelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := functions.ParseFunctionID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsFunctionJavaScriptUDAV0ToV1{},
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
				ValidateFunc: validate.FunctionName,
			},

			"stream_analytics_job_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: streamingjobs.ValidateStreamingJobID,
			},

			"input": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"any",
								"array",
								"bigint",
								"datetime",
								"float",
								"nvarchar(max)",
								"record",
							}, false),
						},

						"configuration_parameter": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"output": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"any",
								"array",
								"bigint",
								"datetime",
								"float",
								"nvarchar(max)",
								"record",
							}, false),
						},
					},
				},
			},

			"script": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceStreamAnalyticsFunctionUDACreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	jobId, err := streamingjobs.ParseStreamingJobID(d.Get("stream_analytics_job_id").(string))
	if err != nil {
		return err
	}

	id := functions.NewFunctionID(subscriptionId, jobId.ResourceGroupName, jobId.StreamingJobName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_stream_analytics_function_javascript_uda", id.ID())
	}

	props := functions.Function{
		Properties: &functions.AggregateFunctionProperties{
			Properties: &functions.FunctionConfiguration{
				Binding: &functions.JavaScriptFunctionBinding{
					Properties: &functions.JavaScriptFunctionBindingProperties{
						Script: utils.String(d.Get("script").(string)),
					},
				},
				Inputs: expandStreamAnalyticsFunctionUDAInputs(d.Get("input").([]interface{})),
				Output: expandStreamAnalyticsFunctionUDAOutput(d.Get("output").([]interface{})),
			},
		},
	}

	var opts functions.CreateOrReplaceOperationOptions
	if _, err := client.CreateOrReplace(ctx, id, props, opts); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStreamAnalyticsFunctionUDARead(d, meta)
}

func resourceStreamAnalyticsFunctionUDARead(d *pluginsdk.ResourceData, meta interface{}) error {
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

			script := ""
			if v := binding.Properties.Script; v != nil {
				script = *v
			}
			d.Set("script", script)

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

func resourceStreamAnalyticsFunctionUDAUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := functions.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	props := functions.Function{
		Properties: &functions.AggregateFunctionProperties{
			Properties: &functions.FunctionConfiguration{
				Binding: &functions.JavaScriptFunctionBinding{
					Properties: &functions.JavaScriptFunctionBindingProperties{
						Script: utils.String(d.Get("script").(string)),
					},
				},
				Inputs: expandStreamAnalyticsFunctionUDAInputs(d.Get("input").([]interface{})),
				Output: expandStreamAnalyticsFunctionUDAOutput(d.Get("output").([]interface{})),
			},
		},
	}

	var opts functions.UpdateOperationOptions
	if _, err := client.Update(ctx, *id, props, opts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsFunctionUDARead(d, meta)
}

func resourceStreamAnalyticsFunctionUDADelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := functions.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandStreamAnalyticsFunctionUDAInputs(input []interface{}) *[]functions.FunctionInput {
	outputs := make([]functions.FunctionInput, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})
		variableType := v["type"].(string)
		outputs = append(outputs, functions.FunctionInput{
			DataType:                 utils.String(variableType),
			IsConfigurationParameter: utils.Bool(v["configuration_parameter"].(bool)),
		})
	}

	return &outputs
}

func flattenStreamAnalyticsFunctionUDAInputs(input *[]functions.FunctionInput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		var variableType string
		if v.DataType != nil {
			variableType = *v.DataType
		}

		var isConfigurationParameter bool
		if v.IsConfigurationParameter != nil {
			isConfigurationParameter = *v.IsConfigurationParameter
		}

		outputs = append(outputs, map[string]interface{}{
			"type":                    variableType,
			"configuration_parameter": isConfigurationParameter,
		})
	}

	return outputs
}

func expandStreamAnalyticsFunctionUDAOutput(input []interface{}) *functions.FunctionOutput {
	output := input[0].(map[string]interface{})

	dataType := output["type"].(string)
	return &functions.FunctionOutput{
		DataType: utils.String(dataType),
	}
}

func flattenStreamAnalyticsFunctionUDAOutput(input *functions.FunctionOutput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var variableType string
	if input.DataType != nil {
		variableType = *input.DataType
	}

	return []interface{}{
		map[string]interface{}{
			"type": variableType,
		},
	}
}
