package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2020-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
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
			_, err := parse.FunctionID(id)
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
				ValidateFunc: validate.FunctionName,
			},

			"stream_analytics_job_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StreamingJobID,
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

	jobId, err := parse.StreamingJobID(d.Get("stream_analytics_job_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFunctionID(subscriptionId, jobId.ResourceGroup, jobId.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_stream_analytics_function_javascript_uda", id.ID())
	}

	props := streamanalytics.Function{
		Properties: &streamanalytics.AggregateFunctionProperties{
			Type: streamanalytics.TypeBasicFunctionPropertiesTypeAggregate,
			FunctionConfiguration: &streamanalytics.FunctionConfiguration{
				Binding: &streamanalytics.JavaScriptFunctionBinding{
					Type: streamanalytics.TypeBasicFunctionBindingTypeMicrosoftStreamAnalyticsJavascriptUdf,
					JavaScriptFunctionBindingProperties: &streamanalytics.JavaScriptFunctionBindingProperties{
						Script: utils.String(d.Get("script").(string)),
					},
				},
				Inputs: expandStreamAnalyticsFunctionUDAInputs(d.Get("input").([]interface{})),
				Output: expandStreamAnalyticsFunctionUDAOutput(d.Get("output").([]interface{})),
			},
		},
	}

	if _, err := client.CreateOrReplace(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, "", ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStreamAnalyticsFunctionUDARead(d, meta)
}

func resourceStreamAnalyticsFunctionUDARead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)

	jobId := parse.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingjobName)
	d.Set("stream_analytics_job_id", jobId.ID())

	if props := resp.Properties; props != nil {
		aggregateProps, ok := props.AsAggregateFunctionProperties()
		if !ok {
			return fmt.Errorf("converting Props to a Aggregate Function")
		}

		binding, ok := aggregateProps.Binding.AsJavaScriptFunctionBinding()
		if !ok {
			return fmt.Errorf("converting Binding to a JavaScript Function Binding")
		}

		if bindingProps := binding.JavaScriptFunctionBindingProperties; bindingProps != nil {
			d.Set("script", bindingProps.Script)
		}

		if err := d.Set("input", flattenStreamAnalyticsFunctionUDAInputs(aggregateProps.Inputs)); err != nil {
			return fmt.Errorf("flattening `input`: %+v", err)
		}

		if err := d.Set("output", flattenStreamAnalyticsFunctionUDAOutput(aggregateProps.Output)); err != nil {
			return fmt.Errorf("flattening `output`: %+v", err)
		}
	}

	return nil
}

func resourceStreamAnalyticsFunctionUDAUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionID(d.Id())
	if err != nil {
		return err
	}

	props := streamanalytics.Function{
		Properties: &streamanalytics.AggregateFunctionProperties{
			Type: streamanalytics.TypeBasicFunctionPropertiesTypeAggregate,
			FunctionConfiguration: &streamanalytics.FunctionConfiguration{
				Binding: &streamanalytics.JavaScriptFunctionBinding{
					Type: streamanalytics.TypeBasicFunctionBindingTypeMicrosoftStreamAnalyticsJavascriptUdf,
					JavaScriptFunctionBindingProperties: &streamanalytics.JavaScriptFunctionBindingProperties{
						Script: utils.String(d.Get("script").(string)),
					},
				},
				Inputs: expandStreamAnalyticsFunctionUDAInputs(d.Get("input").([]interface{})),
				Output: expandStreamAnalyticsFunctionUDAOutput(d.Get("output").([]interface{})),
			},
		},
	}

	if _, err := client.Update(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, ""); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsFunctionUDARead(d, meta)
}

func resourceStreamAnalyticsFunctionUDADelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.Name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandStreamAnalyticsFunctionUDAInputs(input []interface{}) *[]streamanalytics.FunctionInput {
	outputs := make([]streamanalytics.FunctionInput, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})
		variableType := v["type"].(string)
		outputs = append(outputs, streamanalytics.FunctionInput{
			DataType:                 utils.String(variableType),
			IsConfigurationParameter: utils.Bool(v["configuration_parameter"].(bool)),
		})
	}

	return &outputs
}

func flattenStreamAnalyticsFunctionUDAInputs(input *[]streamanalytics.FunctionInput) []interface{} {
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

func expandStreamAnalyticsFunctionUDAOutput(input []interface{}) *streamanalytics.FunctionOutput {
	output := input[0].(map[string]interface{})

	dataType := output["type"].(string)
	return &streamanalytics.FunctionOutput{
		DataType: utils.String(dataType),
	}
}

func flattenStreamAnalyticsFunctionUDAOutput(input *streamanalytics.FunctionOutput) []interface{} {
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
