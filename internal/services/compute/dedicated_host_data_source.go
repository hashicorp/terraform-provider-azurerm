package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDedicatedHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDedicatedHostRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DedicatedHostName(),
			},

			"dedicated_host_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DedicatedHostGroupName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDedicatedHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDedicatedHostID(subscriptionId, d.Get("resource_group_name").(string), d.Get("dedicated_host_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.HostName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("dedicated_host_group_name", id.HostGroupName)

	return tags.FlattenAndSet(d, resp.Tags)
}
