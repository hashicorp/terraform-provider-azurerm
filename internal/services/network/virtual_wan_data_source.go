// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceVirtualWanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVirtualWANID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.VirtualWansGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.VirtualWanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("allow_branch_to_branch_traffic", pointer.From(props.AllowBranchToBranchTraffic))
			d.Set("disable_vpn_encryption", pointer.From(props.DisableVpnEncryption))
			if err := d.Set("office365_local_breakout_category", pointer.From(props.Office365LocalBreakoutCategory)); err != nil {
				return fmt.Errorf("setting `office365_local_breakout_category`: %v", err)
			}
			d.Set("office365_local_breakout_category", pointer.From(props.Office365LocalBreakoutCategory))
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
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func flattenVirtualWanProperties(input *[]virtualwans.SubResource) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	output := make([]interface{}, 0)
	for _, v := range *input {
		if v.Id != nil {
			output = append(output, *v.Id)
		}
	}
	return output
}
