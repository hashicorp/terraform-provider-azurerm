package containers

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceContainerRegistryReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerRegistryReplicationCreate,
		Read:   resourceContainerRegistryReplicationRead,
		Update: resourceContainerRegistryReplicationUpdate,
		Delete: resourceContainerRegistryReplicationDelete,

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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\W_]{5,50}$`), "alpha numeric characters only are allowed and between 5 and 50 characters."),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateContainerRegistryName,
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceContainerRegistryReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ReplicationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	registryName := d.Get("registry_name").(string)
	name := d.Get("name").(string)

	id := parse.NewReplicationID(subscriptionId, resourceGroup, registryName, name)

	existing, err := client.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Container Registry Replication %q: %s", id, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_container_registry_replication", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	replication := containerregistry.Replication{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, registryName, name, replication)
	if err != nil {
		return fmt.Errorf("Error creating Container Registry Replication %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Container Registry Replication %q: %+v", id, err)
	}

	_, err = client.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry Replication %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryReplicationRead(d, meta)
}

func resourceContainerRegistryReplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ReplicationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReplicationID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	replication := containerregistry.ReplicationUpdateParameters{
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.RegistryName, id.Name, replication)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry Replication %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Container Registry Replication %q: %+v", id, err)
	}

	return resourceContainerRegistryReplicationRead(d, meta)
}

func resourceContainerRegistryReplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ReplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry Replication %q was not found in Resource Group %q for Registry %q", id.Name, id.ResourceGroup, id.RegistryName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Registry Replication %q (Resource Group %q, Registry %q): %+v", id.Name, id.ResourceGroup, id.RegistryName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("registry_name", id.RegistryName)
	d.Set("name", id.Name)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceContainerRegistryReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ReplicationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReplicationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry Replication '%s': %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry Replication '%s': %+v", id, err)
	}

	return nil
}
