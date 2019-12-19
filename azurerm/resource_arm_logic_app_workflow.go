package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2016-06-01/logic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var logicAppResourceName = "azurerm_logic_app"

func resourceArmLogicAppWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogicAppWorkflowCreate,
		Read:   resourceArmLogicAppWorkflowRead,
		Update: resourceArmLogicAppWorkflowUpdate,
		Delete: resourceArmLogicAppWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			// TODO: should Parameters be split out into their own object to allow validation on the different sub-types?
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"workflow_schema": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
			},

			"workflow_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "1.0.0.0",
			},

			"tags": tags.Schema(),

			"access_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLogicAppWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Logic App Workflow creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Logic App Workflow %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_logic_app_workflow", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	parameters := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]interface{}))

	workflowSchema := d.Get("workflow_schema").(string)
	workflowVersion := d.Get("workflow_version").(string)
	t := d.Get("tags").(map[string]interface{})

	properties := logic.Workflow{
		Location: utils.String(location),
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: &map[string]interface{}{
				"$schema":        workflowSchema,
				"contentVersion": workflowVersion,
				"actions":        make(map[string]interface{}),
				"triggers":       make(map[string]interface{}),
			},
			Parameters: parameters,
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties); err != nil {
		return fmt.Errorf("[ERROR] Error creating Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read Logic App Workflow %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogicAppWorkflowRead(d, meta)
}

func resourceArmLogicAppWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["workflows"]

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(name, logicAppResourceName)
	defer locks.UnlockByName(name, logicAppResourceName)

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.WorkflowProperties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	parameters := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]interface{}))
	t := d.Get("tags").(map[string]interface{})

	properties := logic.Workflow{
		Location: utils.String(location),
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: read.WorkflowProperties.Definition,
			Parameters: parameters,
		},
		Tags: tags.Expand(t),
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, name, properties); err != nil {
		return fmt.Errorf("Error updating Logic App Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmLogicAppWorkflowRead(d, meta)
}

func resourceArmLogicAppWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workflows"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Logic App Workflow %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkflowProperties; props != nil {
		parameters := flattenLogicAppWorkflowParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}

		d.Set("access_endpoint", props.AccessEndpoint)

		if definition := props.Definition; definition != nil {
			if v, ok := definition.(map[string]interface{}); ok {
				d.Set("workflow_schema", v["$schema"].(string))
				d.Set("workflow_version", v["contentVersion"].(string))
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogicAppWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workflows"]

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(name, logicAppResourceName)
	defer locks.UnlockByName(name, logicAppResourceName)

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing delete request for Logic App Workflow %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandLogicAppWorkflowParameters(input map[string]interface{}) map[string]*logic.WorkflowParameter {
	output := make(map[string]*logic.WorkflowParameter)

	for k, v := range input {
		output[k] = &logic.WorkflowParameter{
			Type:  logic.ParameterTypeString,
			Value: v.(string),
		}
	}

	return output
}

func flattenLogicAppWorkflowParameters(input map[string]*logic.WorkflowParameter) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.Value.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}
