package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistryAgentPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryAgentPoolCreate,
		Read:   resourceContainerRegistryAgentPoolRead,
		Update: resourceContainerRegistryAgentPoolUpdate,
		Delete: resourceContainerRegistryAgentPoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ContainerRegistryAgentPoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 20),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate2.ContainerRegistryName,
			},

			"instance_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1,
			},

			"tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "S1",
				ValidateFunc: validation.StringInSlice([]string{
					"S1",
					"S2",
					"S3",
					"I6",
				}, false),
			},

			"virtual_network_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceContainerRegistryAgentPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryAgentPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Agent Pool creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)

	id := parse.NewContainerRegistryAgentPoolID(subscriptionId, resourceGroup, containerRegistryName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Container Registry Agent Pool %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry_agent_pool", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	count := d.Get("instance_count").(int)
	tier := d.Get("tier").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := containerregistry.AgentPool{
		Location: &location,
		AgentPoolProperties: &containerregistry.AgentPoolProperties{
			// @favoretti: Only Linux is supported
			Os:    containerregistry.OSLinux,
			Count: utils.Int32(int32(count)),
			Tier:  &tier,
		},

		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("virtual_network_subnet_id"); ok {
		parameters.AgentPoolProperties.VirtualNetworkSubnetResourceID = utils.String(v.(string))
	}

	future, err := client.Create(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Container Registry Agent Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Container Registry Agent Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryAgentPoolRead(d, meta)
}

func resourceContainerRegistryAgentPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryAgentPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Agent Pool creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)

	id := parse.NewContainerRegistryAgentPoolID(subscriptionId, resourceGroup, containerRegistryName, name)

	count := d.Get("instance_count").(int)

	parameters := containerregistry.AgentPoolUpdateParameters{
		AgentPoolPropertiesUpdateParameters: &containerregistry.AgentPoolPropertiesUpdateParameters{
			Count: utils.Int32(int32(count)),
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("updating Container Registry Agent Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Container Registry Agent Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryAgentPoolRead(d, meta)
}

func resourceContainerRegistryAgentPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryAgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryAgentPoolID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	registryName := id.RegistryName
	name := id.AgentPoolName

	resp, err := client.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry Agent Pool %q was not found in Container Registry %q", name, registryName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Registry Agent Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("container_registry_name", registryName)

	location := resp.Location
	if location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AgentPoolProperties; props != nil {
		d.Set("instance_count", props.Count)
		d.Set("tier", props.Tier)

		if resp.AgentPoolProperties.VirtualNetworkSubnetResourceID != nil {
			d.Set("virtual_network_subnet_id", props.VirtualNetworkSubnetResourceID)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceContainerRegistryAgentPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryAgentPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	name := id.AgentPoolName
	resourceGroup := id.ResourceGroup
	registryName := id.RegistryName

	future, err := client.Delete(ctx, resourceGroup, registryName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("issuing Azure ARM delete request of Container Registry Agent Pool '%s': %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("issuing Azure ARM delete request of Container Registry Agent Pool '%s': %+v", name, err)
	}

	return nil
}
