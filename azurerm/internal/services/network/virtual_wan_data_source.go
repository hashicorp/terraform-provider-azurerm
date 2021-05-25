package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceVirtualWan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualWanRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"allow_branch_to_branch_traffic": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
			"disable_vpn_encryption": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
			"office365_local_breakout_category": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"virtual_hub_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"vpn_site_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"location": azure.SchemaLocationForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVirtualWanRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Virtual Wan %q (resource group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualWanProperties; props != nil {
		d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		if err := d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory); err != nil {
			return fmt.Errorf("error setting `office365_local_breakout_category`: %v", err)
		}
		d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory)
		if err := d.Set("sku", props.Type); err != nil {
			return fmt.Errorf("error setting `sku`: %v", err)
		}
		d.Set("sku", props.Type)
		if err := d.Set("virtual_hub_ids", flattenVirtualWanProperties(props.VirtualHubs)); err != nil {
			return fmt.Errorf("error setting `virtual_hubs`: %v", err)
		}
		if err := d.Set("vpn_site_ids", flattenVirtualWanProperties(props.VpnSites)); err != nil {
			return fmt.Errorf("error setting `vpn_sites`: %v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenVirtualWanProperties(input *[]network.SubResource) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	output := make([]interface{}, 0)
	for _, v := range *input {
		if v.ID != nil {
			output = append(output, *v.ID)
		}
	}
	return output
}
