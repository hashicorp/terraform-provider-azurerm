package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDatabricksWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDatabricksWorkspaceCreate,
		Read:   resourceArmDatabricksWorkspaceRead,
		Update: resourceArmDatabricksWorkspaceUpdate,
		Delete: resourceArmDatabricksWorkspaceDelete,

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "standard",
				ValidateFunc: validation.StringInSlice([]string{
					"standard",
					"premium",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDatabricksWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricksWorkspacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Databricks Workspace creation.")

	subscriptionID := meta.(*ArmClient).subscriptionId
	name := d.Get("workspace_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	skuName := d.Get("sku").(string)

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	managedResourceGroupID := fmt.Sprintf("/subscriptions/%s/resourceGroups/databricks-rg-%s", subscriptionID, resGroup)

	properties := databricks.WorkspaceProperties{
		ManagedResourceGroupID: &managedResourceGroupID,
	}

	workspace := databricks.Workspace{
		Sku: &databricks.Sku{
			Name: &skuName,
		},
		Location:            &location,
		WorkspaceProperties: &properties,
		Tags:                expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, workspace, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Databricks Workspace Instance %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Databricks Workspace (%s) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    workspaceStateRefreshFunc(ctx, client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Databricks Workspace (%s) to become available: %s", name, err)
	}

	d.SetId(*read.ID)

	return resourceArmDatabricksWorkspaceRead(d, meta)
}

func resourceArmDatabricksWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricksWorkspacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Databricks Workspace update.")

	name := d.Get("workspace_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	workspaceUpdate := databricks.WorkspaceUpdate{
		Tags: expandedTags,
	}

	_, err := client.Update(ctx, workspaceUpdate, resGroup, name)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Databricks Workspace %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Databricks Workspace (%s) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    workspaceStateRefreshFunc(ctx, client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Databricks Workspace (%s) to become available: %s", name, err)
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
		return fmt.Errorf("Error Retrieving Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	// covers if the resource has been deleted outside of TF, but is still in the state
	if utils.ResponseWasNotFound(resp.Response) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Databricks Workspace %s: %s", name, err)
	}

	d.Set("workspace_name", name)
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
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	return nil
}

func workspaceStateRefreshFunc(ctx context.Context, client databricks.WorkspacesClient, resourceGroupName string, workspaceName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, workspaceName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in workspaceStateRefreshFunc to Azure ARM for Databricks Workspace '%s' (RG: '%s'): %s", workspaceName, resourceGroupName, err)
		}

		return res, string(res.ProvisioningState), nil
	}
}
