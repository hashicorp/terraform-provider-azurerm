package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmLogicAppWorkflow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLogicAppWorkflowRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			// TODO: should Parameters be split out into their own object to allow validation on the different sub-types?
			"parameters": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"workflow_schema": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workflow_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),

			"access_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceArmLogicAppWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logicWorkflowsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Logic App Workflow %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.WorkflowProperties; props != nil {
		parameters := flattenLogicAppDataSourceWorkflowParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}

		d.Set("access_endpoint", props.AccessEndpoint)

		if definition := props.Definition; definition != nil {
			if v, ok := definition.(map[string]interface{}); ok {
				schema := v["$schema"].(string)
				version := v["contentVersion"].(string)
				d.Set("workflow_schema", schema)
				d.Set("workflow_version", version)
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenLogicAppDataSourceWorkflowParameters(input map[string]*logic.WorkflowParameter) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			output[k] = v.Value.(string)
		}
	}

	return output
}
