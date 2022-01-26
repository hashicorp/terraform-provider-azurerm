package containers

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceContainerRegistryScopeMap() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceContainerRegistryScopeMapRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceContainerRegistryScopeMapRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ScopeMapsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewContainerRegistryScopeMapID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.ScopeMapName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("description", resp.Description)
	d.Set("actions", utils.FlattenStringSlice(resp.Actions))

	return nil
}
