package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDatabricksWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDatabricksWorkspaceCreateUpdate,
		Read:   resourceArmDatabricksWorkspaceRead,
		Update: resourceArmDatabricksWorkspaceCreateUpdate,
		Delete: resourceArmDatabricksWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatabricksWorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"standard",
					"premium",
				}, false),
			},

			"tags": tagsSchema(),

			"managed_resource_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"managed_resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDatabricksWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricks.WorkspacesClient
	ctx := meta.(*ArmClient).StopContext
	subscriptionID := meta.(*ArmClient).subscriptionId

	log.Printf("[INFO] preparing arguments for Azure ARM Databricks Workspace creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Databricks Workspace %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_databricks_workspace", *existing.ID)
		}
	}

	skuName := d.Get("sku").(string)
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)
	var managedResourceGroupID string

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	if managedResourceGroupName == "" {
		//no managed resource group name was provided, we use the default pattern
		log.Printf("[DEBUG][azurerm_databricks_workspace] no managed resource group id was provided, we use the default pattern.")
		managedResourceGroupID = fmt.Sprintf("/subscriptions/%s/resourceGroups/databricks-rg-%s", subscriptionID, resourceGroup)
	} else {
		log.Printf("[DEBUG][azurerm_databricks_workspace] a managed group name was provided: %q", managedResourceGroupName)
		managedResourceGroupID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, managedResourceGroupName)
	}

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

	future, err := client.CreateOrUpdate(ctx, workspace, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error creating/updating Databricks Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of Databricks Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Databricks Workspace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Databricks Workspace %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDatabricksWorkspaceRead(d, meta)
}

func resourceArmDatabricksWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricks.WorkspacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Databricks Workspace %q was not found in Resource Group %q - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Databricks Workspace %s: %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}

	if props := resp.WorkspaceProperties; props != nil {
		managedResourceGroupID, err := azure.ParseAzureResourceID(*props.ManagedResourceGroupID)
		if err != nil {
			return err
		}
		d.Set("managed_resource_group_id", props.ManagedResourceGroupID)
		d.Set("managed_resource_group_name", managedResourceGroupID.ResourceGroup)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDatabricksWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).databricks.WorkspacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Databricks Workspace %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	return nil
}

func validateDatabricksWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q type to be string", k))
		return warnings, errors
	}

	// Cannot be empty
	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return warnings, errors
	}

	// First, second, and last characters must be a letter or number with a total length between 3 to 64 characters
	// NOTE: Restricted name to 30 characters because that is the restriction in Azure Portal even though the API supports 64 characters
	if !regexp.MustCompile("^[a-zA-Z0-9]{2}[-a-zA-Z0-9]{0,27}[a-zA-Z0-9]{1}$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be 3 - 30 characters in length", k))
		errors = append(errors, fmt.Errorf("%q first, second, and last characters must be a letter or number", k))
		errors = append(errors, fmt.Errorf("%q can only contain letters, numbers, and hyphens", k))
	}

	// No consecutive hyphens
	if regexp.MustCompile("(--)").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must not contain any consecutive hyphens", k))
	}

	return warnings, errors
}
