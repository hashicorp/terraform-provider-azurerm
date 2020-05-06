package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIpPrefix() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIpPrefixRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"prefix_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ip_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaZonesComputed(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPublicIpPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Public IP prefix %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("Error retrieving Public IP Prefix %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
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
