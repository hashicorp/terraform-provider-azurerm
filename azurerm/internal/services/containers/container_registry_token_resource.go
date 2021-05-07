package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2020-11-01-preview/containerregistry"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceContainerRegistryToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerRegistryTokenCreate,
		Read:   resourceContainerRegistryTokenRead,
		Update: resourceContainerRegistryTokenUpdate,
		Delete: resourceContainerRegistryTokenDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryTokenName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"container_registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},

			"scope_map_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryScopeMapID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceContainerRegistryTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry token creation.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing token %q in Container Registry %q (Resource Group %q): %s", name, containerRegistryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry_token", *existing.ID)
		}
	}

	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := containerregistry.TokenStatusEnabled

	if !enabled {
		status = containerregistry.TokenStatusDisabled
	}

	parameters := containerregistry.Token{
		TokenProperties: &containerregistry.TokenProperties{
			ScopeMapID: utils.String(scopeMapID),
			Status:     status,
		},
	}

	future, err := client.Create(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating token %q in Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving token %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read token %q for Container Registry %q (resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceContainerRegistryTokenRead(d, meta)
}
func resourceContainerRegistryTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry token update.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)
	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := containerregistry.TokenStatusEnabled

	if !enabled {
		status = containerregistry.TokenStatusDisabled
	}

	parameters := containerregistry.TokenUpdateParameters{
		TokenUpdateProperties: &containerregistry.TokenUpdateProperties{
			ScopeMapID: utils.String(scopeMapID),
			Status:     status,
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating token %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read token %q (Container Registry %q, resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceContainerRegistryTokenRead(d, meta)
}

func resourceContainerRegistryTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryTokenID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Token %q was not found in Container Registry %q in Resource Group %q", id.TokenName, id.RegistryName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on token %q in Azure Container Registry %q (Resource Group %q): %+v", id.TokenName, id.RegistryName, id.ResourceGroup, err)
	}

	status := true
	if resp.Status == containerregistry.TokenStatusDisabled {
		status = false
	}

	d.Set("name", resp.Name)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("enabled", status)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("scope_map_id", resp.ScopeMapID)

	return nil
}

func resourceContainerRegistryTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryTokenID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry token '%s': %+v", id.TokenName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry token '%s': %+v", id.TokenName, err)
	}

	return nil
}
