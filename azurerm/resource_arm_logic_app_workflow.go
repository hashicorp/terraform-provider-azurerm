package azurerm

import (
	"fmt"
	"log"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogicAppWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogicAppWorkflowCreateUpdate,
		Read:   resourceArmLogicAppWorkflowRead,
		Update: resourceArmLogicAppWorkflowCreateUpdate,
		Delete: resourceArmLogicAppWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"trigger": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recurrence": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"frequency": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Month",
											"Week",
											"Day",
											"Hour",
											"Minute",
											"Hour",
											"Second",
										}, true),
									},
									"interval": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"method": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"function": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"function_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"body": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),

			"access_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLogicAppWorkflowCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Logic App Workflow creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	definition := expandLogicAppWorkflowDefinition(d)
	parameters := expandLogicAppWorkflowParameters(d)
	tags := d.Get("tags").(map[string]interface{})

	properties := logic.Workflow{
		Location: utils.String(location),
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: parameters,
		},
		Tags: expandTags(tags),
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read Logic App Workflow %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogicAppWorkflowRead(d, meta)
}

func resourceArmLogicAppWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workflows"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.WorkflowProperties; props != nil {
		parameters := flattenLogicAppWorkflowParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error flattening `parameters`: %+v", err)
		}

		triggers, err := flattenLogicAppWorkflowTriggers(props.Definition)
		if err != nil {
			return fmt.Errorf("Error flattening `trigger`: %+v", err)
		}
		if err := d.Set("trigger", triggers); err != nil {
			return fmt.Errorf("Error setting `trigger`: %+v", err)
		}

		actions, err := flattenLogicAppWorkflowActions(props.Definition)
		if err != nil {
			return fmt.Errorf("Error flattening `action`: %+v", err)
		}
		if err := d.Set("action", actions); err != nil {
			return fmt.Errorf("Error setting `action`: %+v", err)
		}

		d.Set("access_endpoint", props.AccessEndpoint)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmLogicAppWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workflows"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing delete request for Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandLogicAppWorkflowDefinition(d *schema.ResourceData) interface{} {
	triggers := expandLogicAppWorkflowTriggers(d)
	actions := expandLogicAppWorkflowActions(d)

	definition := map[string]interface{}{
		// these could potentially be user-configurable in future?
		"$schema":        "http://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json",
		"contentVersion": "1.0.0.0",
		"actions":        actions,
		"triggers":       triggers,
	}
	return definition
}

func expandLogicAppWorkflowActions(d *schema.ResourceData) interface{} {
	output := make(map[string]interface{}, 0)
	input := d.Get("action").([]interface{})

	for _, trigger := range input {
		triggerVal := trigger.(map[string]interface{})

		https := triggerVal["http"].([]interface{})
		for _, http := range https {
			httpVal := http.(map[string]interface{})

			name := httpVal["name"].(string)
			method := httpVal["method"].(string)
			uri := httpVal["uri"].(string)

			output[name] = map[string]interface{}{
				"type": "http",
				"inputs": map[string]interface{}{
					"method": method,
					"uri":    uri,
				},
			}
		}

		functions := triggerVal["function"].([]interface{})
		for _, function := range functions {
			functionVal := function.(map[string]interface{})

			name := functionVal["name"].(string)
			functionId := functionVal["function_id"].(string)
			body := functionVal["body"].(string)

			output[name] = map[string]interface{}{
				"type": "function",
				"inputs": map[string]interface{}{
					"body": body,
					"function": map[string]interface{}{
						"id": functionId,
					},
				},
			}
		}
	}

	return output
}

func flattenLogicAppWorkflowActions(input interface{}) ([]interface{}, error) {
	definition, ok := input.(map[string]interface{})
	if !ok {
		return []interface{}{}, nil
	}

	actionsVal, ok := definition["actions"]
	if !ok {
		return []interface{}{}, nil
	}

	httpOutputs := make([]interface{}, 0)
	functionOutputs := make([]interface{}, 0)

	actions := actionsVal.(map[string]interface{})
	for name, val := range actions {
		v := val.(map[string]interface{})
		actionType := v["type"].(string)

		if strings.EqualFold(actionType, "http") {
			inputs := v["inputs"].(map[string]interface{})
			method := inputs["method"].(string)
			uri := inputs["uri"].(string)
			output := map[string]interface{}{
				"name":   name,
				"method": method,
				"uri":    uri,
			}
			httpOutputs = append(httpOutputs, output)
		} else if strings.EqualFold(actionType, "function") {
			inputs := v["input"].(map[string]interface{})
			body := inputs["body"].(string)
			function := inputs["function"].(map[string]interface{})
			functionId := function["id"].(string)
			output := map[string]interface{}{
				"name":        name,
				"function_id": functionId,
				"body":        body,
			}
			functionOutputs = append(functionOutputs, output)
		} else {
			return nil, fmt.Errorf("Unsupported Action Type %q", actionType)
		}
	}

	output := map[string]interface{}{
		"http":     httpOutputs,
		"function": functionOutputs,
	}

	return []interface{}{output}, nil
}

func expandLogicAppWorkflowTriggers(d *schema.ResourceData) interface{} {
	output := make(map[string]interface{}, 0)

	input := d.Get("trigger").([]interface{})
	for _, triggerVal := range input {
		val := triggerVal.(map[string]interface{})

		recurrenceInput := val["recurrence"].([]interface{})
		for _, recurrenceVal := range recurrenceInput {
			v := recurrenceVal.(map[string]interface{})
			name := v["name"].(string)
			frequency := v["frequency"].(string)
			interval := v["interval"].(int)
			output[name] = map[string]interface{}{
				"type": "recurrence",
				"recurrence": map[string]interface{}{
					"frequency": frequency,
					"interval":  interval,
				},
			}
		}
	}

	return output
}

func flattenLogicAppWorkflowTriggers(input interface{}) ([]interface{}, error) {
	definition, ok := input.(map[string]interface{})
	if !ok {
		return []interface{}{}, nil
	}

	triggersVal, ok := definition["triggers"]
	if !ok {
		return []interface{}{}, nil
	}

	recurrenceOutputs := make([]interface{}, 0)

	triggers := triggersVal.(map[string]interface{})
	for name, val := range triggers {
		v := val.(map[string]interface{})
		triggerType := v["type"].(string)

		if strings.EqualFold(triggerType, "recurrence") {
			recurrence := v["recurrence"].(map[string]interface{})

			frequency := recurrence["frequency"].(string)
			interval := int(recurrence["interval"].(float64))
			output := map[string]interface{}{
				"name":      name,
				"frequency": frequency,
				"interval":  interval,
			}
			recurrenceOutputs = append(recurrenceOutputs, output)
		} else {
			return nil, fmt.Errorf("Unsupported Trigger Type %q", triggerType)
		}
	}

	output := map[string]interface{}{
		"recurrence": recurrenceOutputs,
	}

	return []interface{}{output}, nil
}

func expandLogicAppWorkflowParameters(d *schema.ResourceData) map[string]*logic.WorkflowParameter {
	output := make(map[string]*logic.WorkflowParameter, 0)
	input := d.Get("parameters").(map[string]interface{})

	for k, v := range input {
		output[k] = &logic.WorkflowParameter{
			Type:  logic.ParameterTypeString,
			Value: v.(string),
		}
	}

	return output
}

func flattenLogicAppWorkflowParameters(input map[string]*logic.WorkflowParameter) map[string]interface{} {
	output := make(map[string]interface{}, 0)

	for k, v := range input {
		if v != nil {
			output[k] = v.Value.(string)
		}
	}

	return output
}
