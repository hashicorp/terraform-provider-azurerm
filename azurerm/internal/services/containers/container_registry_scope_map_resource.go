package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerRegistryScopeMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerRegistryScopeMapCreate,
		Read:   resourceArmContainerRegistryScopeMapRead,
		Update: resourceArmContainerRegistryScopeMapUpdate,
		Delete: resourceArmContainerRegistryScopeMapDelete,
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
			"description": {
				Type: schema.TypeString,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"container_registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRMContainerRegistryName,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: validation.StringInSlice([]string{
					string(containerregistry.TokenStatusDisabled),
					string(containerregistry.TokenStatusEnabled),
				}, false),
			},
		},
	}
}

func resourceArmContainerRegistryScopeMapCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry scope map creation.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing scope map %q in Container Registry %q (Resource Group %q): %s", name, containerRegistryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry_scope_map", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	actions := d.Get("actions").([]interface{})

	parameters := containerregistry.ScopeMap{
		ScopeMapProperties: &containerregistry.ScopeMapProperties{
			Description: utils.String(description),
			Actions:     utils.ExpandStringSlice(actions),
		},
	}

	future, err := client.Create(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating scope map %q in Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of scope map %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving scope map %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read scope map %q for Container Registry %q (resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryScopeMapRead(d, meta)
}

func resourceArmContainerRegistryScopeMapUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry scope map update.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	actions := d.Get("actions").([]interface{})

	parameters := containerregistry.ScopeMapUpdateParameters{
		ScopeMapPropertiesUpdateParameters: &containerregistry.ScopeMapPropertiesUpdateParameters{
			Description: utils.String(description),
			Actions:     utils.ExpandStringSlice(actions),
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating scope map %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of scope map %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving scope map %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read scope map %q (Container Registry %q, resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryTokenRead(d, meta)
}

func resourceArmContainerRegistryScopeMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	containerRegistryName := id.Path["container_registry_name"]
	name := id.Path["name"]

	resp, err := client.Get(ctx, resourceGroup, containerRegistryName, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Scope Map %q was not found in Container Registry %q in Resource Group %q", name, containerRegistryName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on scope map %q in Azure Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("actions", utils.FlattenStringSlice(resp.Actions))

	return nil
}

func resourceArmContainerRegistryScopeMapDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	containerRegistryName := id.Path["container_registry_name"]
	name := id.Path["name"]

	future, err := client.Delete(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry scope map '%s': %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry scope map '%s': %+v", name, err)
	}

	return nil
}
