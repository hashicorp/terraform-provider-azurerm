package containers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceContainerRegistryToken() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceContainerRegistryTokenRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryTokenName,
			},

			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"scope_map_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceContainerRegistryTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewContainerRegistryTokenID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	status := true
	if resp.Status == containerregistry.TokenStatusDisabled {
		status = false
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("scope_map_id", resp.ScopeMapID)
	d.Set("enabled", status)

	return nil
}
