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

			"allow_branch_to_branch_traffic": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_vnet_to_vnet_traffic": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_vpn_encryption": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"office365_local_breakout_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_hubs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_sites": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// "address_prefix": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "data": {
			// 	Type: schema.TypeMap,
			// 	Computed: true,
			// }
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
	if props := resp.VirtualWanProperties; props != nil {
		if props.AllowBranchToBranchTraffic != nil {
			d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		}
		if props.AllowVnetToVnetTraffic != nil {
			d.Set("allow_vnet_to_vnet_traffic", props.AllowVnetToVnetTraffic)
		}
		if props.DisableVpnEncryption != nil {
			d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		}
		//bool
		d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory) //string
		d.Set("sku", props.Type)                                                         // string
		d.Set("virtual_hubs", props.VirtualHubs)                                         //list
		d.Set("vpn_sites", props.VpnSites)                                               //list
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
