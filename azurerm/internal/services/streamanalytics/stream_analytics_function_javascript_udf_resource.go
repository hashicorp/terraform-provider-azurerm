package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStreamAnalyticsFunctionUDF() *schema.Resource {
	return &schema.Resource{
		Create: resourceStreamAnalyticsFunctionUDFCreateUpdate,
		Read:   resourceStreamAnalyticsFunctionUDFRead,
		Update: resourceStreamAnalyticsFunctionUDFCreateUpdate,
		Delete: resourceStreamAnalyticsFunctionUDFDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FunctionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"stream_analytics_job_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"input": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"any",
								"datetime",
								"array",
								"bigint",
								"float",
								"nvarchar(max)",
								"record",
							}, false),
						},
					},
				},
			},

			"output": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"any",
								"datetime",
								"array",
								"bigint",
								"float",
								"nvarchar(max)",
								"record",
							}, false),
						},
					},
				},
			},

			"script": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				// TODO: JS diff suppress func?!
			},
		},
	}
}

func resourceStreamAnalyticsFunctionUDFCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Function Javascript UDF creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Function Javascript UDF %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_function_javascript_udf", *existing.ID)
		}
	}

	script := d.Get("script").(string)
	inputsRaw := d.Get("input").([]interface{})
	inputs := expandStreamAnalyticsFunctionInputs(inputsRaw)

	outputRaw := d.Get("output").([]interface{})
	output := expandStreamAnalyticsFunctionOutput(outputRaw)

	function := streamanalytics.Function{
		Properties: &streamanalytics.ScalarFunctionProperties{
			Type: streamanalytics.TypeScalar,
			ScalarFunctionConfiguration: &streamanalytics.ScalarFunctionConfiguration{
				Binding: &streamanalytics.JavaScriptFunctionBinding{
					Type: streamanalytics.TypeMicrosoftStreamAnalyticsJavascriptUdf,
					JavaScriptFunctionBindingProperties: &streamanalytics.JavaScriptFunctionBindingProperties{
						Script: utils.String(script),
					},
				},
				Inputs: inputs,
				Output: output,
			},
		},
	}

	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, function, resourceGroup, jobName, name, "", ""); err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Function Javascript UDF %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Function Javascript UDF %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Function Javascript UDF %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else if _, err := client.Update(ctx, function, resourceGroup, jobName, name, ""); err != nil {
		return fmt.Errorf("Error Updating Stream Analytics Function Javascript UDF %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	return resourceStreamAnalyticsFunctionUDFRead(d, meta)
}

func resourceStreamAnalyticsFunctionUDFRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("stream_analytics_job_name", id.StreamingjobName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		scalarProps, ok := props.AsScalarFunctionProperties()
		if !ok {
			return fmt.Errorf("converting Props to a Scalar Function")
		}

		binding, ok := scalarProps.Binding.AsJavaScriptFunctionBinding()
		if !ok {
			return fmt.Errorf("converting Binding to a JavaScript Function Binding")
		}

		if bindingProps := binding.JavaScriptFunctionBindingProperties; bindingProps != nil {
			d.Set("script", bindingProps.Script)
		}

		if err := d.Set("input", flattenStreamAnalyticsFunctionInputs(scalarProps.Inputs)); err != nil {
			return fmt.Errorf("flattening `input`: %+v", err)
		}

		if err := d.Set("output", flattenStreamAnalyticsFunctionOutput(scalarProps.Output)); err != nil {
			return fmt.Errorf("flattening `output`: %+v", err)
		}
	}

	return nil
}

func resourceStreamAnalyticsFunctionUDFDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.FunctionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FunctionID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.Name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func expandStreamAnalyticsFunctionInputs(input []interface{}) *[]streamanalytics.FunctionInput {
	outputs := make([]streamanalytics.FunctionInput, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})
		variableType := v["type"].(string)
		outputs = append(outputs, streamanalytics.FunctionInput{
			DataType: utils.String(variableType),
		})
	}

	return &outputs
}

func flattenStreamAnalyticsFunctionInputs(input *[]streamanalytics.FunctionInput) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		var variableType string
		if v.DataType != nil {
			variableType = *v.DataType
		}

		outputs = append(outputs, map[string]interface{}{
			"type": variableType,
		})
	}

	return outputs
}

func expandStreamAnalyticsFunctionOutput(input []interface{}) *streamanalytics.FunctionOutput {
	output := input[0].(map[string]interface{})

	dataType := output["type"].(string)
	return &streamanalytics.FunctionOutput{
		DataType: utils.String(dataType),
	}
}

func flattenStreamAnalyticsFunctionOutput(input *streamanalytics.FunctionOutput) []interface{} {
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
