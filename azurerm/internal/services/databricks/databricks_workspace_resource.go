package databricks

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDatabricksWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDatabricksWorkspaceCreateUpdate,
		Read:   resourceArmDatabricksWorkspaceRead,
		Update: resourceArmDatabricksWorkspaceCreateUpdate,
		Delete: resourceArmDatabricksWorkspaceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DatabricksWorkspaceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateDatabricksWorkspaceName,
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
					"trial",
				}, false),
			},

			"tags": tags.Schema(),

			"managed_resource_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"custom_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"no_public_ip": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
						},

						"public_subnet_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"private_subnet_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},

						"virtual_network_id": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},
					},
				},
			},

			"managed_resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDatabricksWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[INFO] preparing arguments for Azure ARM Databricks Workspace creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	if managedResourceGroupName == "" {
		// no managed resource group name was provided, we use the default pattern
		log.Printf("[DEBUG][azurerm_databricks_workspace] no managed resource group id was provided, we use the default pattern.")
		managedResourceGroupID = fmt.Sprintf("/subscriptions/%s/resourceGroups/databricks-rg-%s", subscriptionID, resourceGroup)
	} else {
		log.Printf("[DEBUG][azurerm_databricks_workspace] a managed group name was provided: %q", managedResourceGroupName)
		managedResourceGroupID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, managedResourceGroupName)
	}

	customParamsRaw := d.Get("custom_parameters").([]interface{})
	customParams := expandWorkspaceCustomParameters(customParamsRaw)

	workspace := databricks.Workspace{
		Sku: &databricks.Sku{
			Name: utils.String(skuName),
		},
		Location: utils.String(location),
		WorkspaceProperties: &databricks.WorkspaceProperties{
			ManagedResourceGroupID: &managedResourceGroupID,
			Parameters:             customParams,
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
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabricksWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Databricks Workspace %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Databricks Workspace %s: %s", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

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

		if err := d.Set("custom_parameters", flattenWorkspaceCustomParameters(props.Parameters)); err != nil {
			return fmt.Errorf("Error setting `custom_parameters`: %+v", err)
		}

		d.Set("workspace_url", props.WorkspaceURL)
		d.Set("workspace_id", props.WorkspaceID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDatabricksWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabricksWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Databricks Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Databricks Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func flattenWorkspaceCustomParameters(p *databricks.WorkspaceCustomParameters) []interface{} {
	if p == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if v := p.EnableNoPublicIP; v != nil {
		if v.Value != nil {
			parameters["no_public_ip"] = *v.Value
		}
	}

	if v := p.CustomPrivateSubnetName; v != nil {
		if v.Value != nil {
			parameters["private_subnet_name"] = *v.Value
		}
	}

	if v := p.CustomPublicSubnetName; v != nil {
		if v.Value != nil {
			parameters["public_subnet_name"] = *v.Value
		}
	}

	if v := p.CustomVirtualNetworkID; v != nil {
		if v.Value != nil {
			parameters["virtual_network_id"] = *v.Value
		}
	}

	return []interface{}{parameters}
}

func expandWorkspaceCustomParameters(input []interface{}) *databricks.WorkspaceCustomParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	parameters := databricks.WorkspaceCustomParameters{}

	if v, ok := config["no_public_ip"].(bool); ok {
		parameters.EnableNoPublicIP = &databricks.WorkspaceCustomBooleanParameter{
			Value: &v,
		}
	}

	if v := config["public_subnet_name"].(string); v != "" {
		parameters.CustomPublicSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v := config["private_subnet_name"].(string); v != "" {
		parameters.CustomPrivateSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v := config["virtual_network_id"].(string); v != "" {
		parameters.CustomVirtualNetworkID = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	return &parameters
}

func ValidateDatabricksWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q type to be string", k))
		return warnings, errors
	}

	// The Azure Portal shows the following validation criteria:

	// 1) Cannot be empty
	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		// Treating this as a special case and returning early to match Azure Portal behavior.
		return warnings, errors
	}

	// 2) Must be at least 3 characters:
	if len(v) < 3 {
		errors = append(errors, fmt.Errorf("%q must be at least 3 characters: %q", k, v))
	}

	// 3) The value must have a length of at most 30.
	// NOTE: Restricted name to 30 characters because that is the restriction in Azure Portal even though the API supports 64 characters
	if len(v) > 30 {
		errors = append(errors, fmt.Errorf("%q must be no more than 30 characters: %q", k, v))
	}

	// 4) Only alphanumeric characters, underscores, and hyphens are allowed.
	if !regexp.MustCompile("^[a-zA-Z0-9_-]*$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can contain only alphanumeric characters, underscores, and hyphens: %q", k, v))
	}

	return warnings, errors
}
