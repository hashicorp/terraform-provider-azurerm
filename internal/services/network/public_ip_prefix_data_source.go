package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourcePublicIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePublicIpPrefixRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"prefix_length": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"ip_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaZonesComputed(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourcePublicIpPrefixRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPublicIpPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.PublicIPPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}
	if props := resp.PublicIPPrefixPropertiesFormat; props != nil {
		d.Set("prefix_length", props.PrefixLength)
		d.Set("ip_prefix", props.IPPrefix)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
