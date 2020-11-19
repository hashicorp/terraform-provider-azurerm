package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"allow_branch_to_branch_traffic": {
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
				Optional: true,
				Computed: true,
			},
			"vpn_sites": {
				Type:     schema.TypeString,
				Optional: true,
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

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Virtual Wan %q (resource group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	b := false

	if props := resp.VirtualWanProperties; props != nil {
		log.Printf("[DEBUG] ETIENNE ETIENNE ETIENNE %+v", props)
		if abtbt := *props.AllowBranchToBranchTraffic; props.AllowBranchToBranchTraffic != nil {
			if err := d.Set("allow_branch_to_branch_traffic", abtbt); err != nil {
				return fmt.Errorf("error setting `allow_branch_to_branch_traffic`: %v", err)
			}
		} else {
			if err := d.Set("allow_branch_to_branch_traffic", &b); err != nil {
				return fmt.Errorf("error setting `allow_branch_to_branch_traffic`: %v", err)
			}
		d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		if dve := *props.DisableVpnEncryption; props.DisableVpnEncryption != nil {
			if err := d.Set("disable_vpn_encryption", dve); err != nil {
				return fmt.Errorf("error setting `disable_vpn_encryption`: %v", err)
			}
		} else {
			if err := d.Set("disable_vpn_encryption", &b); err != nil {
				return fmt.Errorf("error setting `disable_vpn_encryption`: %v", err)
			}
		d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		if err := d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory); err != nil {
			return fmt.Errorf("error setting `office365_local_breakout_category`: %v", err)
		}
		if err := d.Set("sku", props.Type); err != nil {
			return fmt.Errorf("error setting `sku`: %v", err)
		}
		if err := d.Set("virtual_hubs", flattenVirtualWanProperties(props.VirtualHubs)); err != nil {
			return fmt.Errorf("error setting `virtual_hubs`: %v", err)
		}
		if err := d.Set("vpn_sites", flattenVirtualWanProperties(props.VpnSites)); err != nil {
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
