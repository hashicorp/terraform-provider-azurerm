// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			"location": commonschema.LocationComputed(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVirtualWanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualWanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.VirtualWanProperties; props != nil {
		d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		if err := d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory); err != nil {
			return fmt.Errorf("setting `office365_local_breakout_category`: %v", err)
		}
		d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory)
		if err := d.Set("sku", props.Type); err != nil {
			return fmt.Errorf("setting `sku`: %v", err)
		}
		d.Set("sku", props.Type)
		if err := d.Set("virtual_hub_ids", flattenVirtualWanProperties(props.VirtualHubs)); err != nil {
			return fmt.Errorf("setting `virtual_hubs`: %v", err)
		}
		if err := d.Set("vpn_site_ids", flattenVirtualWanProperties(props.VpnSites)); err != nil {
			return fmt.Errorf("setting `vpn_sites`: %v", err)
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
