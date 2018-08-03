package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
)

func resourceArmDatabricksWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDatabricksWorkspaceCreateUpdate,
		Read:   resourceArmDatabricksWorkspaceRead,
		Update: resourceArmDatabricksWorkspaceCreateUpdate,
		Delete: resourceArmDatabricksWorkspaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"managed_resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, false),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDatabricksWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricksWorkspacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Databricks Workspace creation.")

	subscriptionID := meta.(*ArmClient).subscriptionId
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	skuName := d.Get("sku").(string)

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	managedResourceGroupID := fmt.Sprintf("/subscriptions/%s/resourceGroups/databricks-rg-%s", subscriptionID, resGroup)

	workspace := databricks.Workspace{
		Sku: &databricks.Sku{
			Name: utils.String(skuName),
		},
		Location: utils.String(location),
		WorkspaceProperties: &databricks.WorkspaceProperties{
			ManagedResourceGroupID: &managedResourceGroupID,
		},
		Tags: expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, workspace, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error creating Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Databricks Workspace %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDatabricksWorkspaceRead(d, meta)
}

func resourceArmDatabricksWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricksWorkspacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Get(ctx, resGroup, name)

	if err != nil {
		// covers if the resource has been deleted outside of TF, but is still in the state
		d.SetId("")
		return fmt.Errorf("Error Retrieving Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Databricks Workspace %s: %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("managed_resource_group_id", props.ManagedResourceGroupID)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDatabricksWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricksWorkspacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
