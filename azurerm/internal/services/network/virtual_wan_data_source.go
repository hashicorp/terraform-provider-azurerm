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

func dataSourceArmVirtualWan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualWanRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_wan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmVirtualWanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Wan %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Virtual Wan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	// if props := resp.ID; props != nil {
	// 	if props.ID != nil {
	// 		d.Set("virtual_wan_id", props.ID)
	// 	}
	// }
	return tags.FlattenAndSet(d, resp.Tags)
}
