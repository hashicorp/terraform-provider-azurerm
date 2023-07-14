// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceNatGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNatGatewayRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.NatGatewayName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"public_ip_address_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_prefix_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zones": commonschema.ZonesMultipleComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceNatGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNatGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	if props := resp.NatGatewayPropertiesFormat; props != nil {
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		d.Set("resource_guid", props.ResourceGUID)

		if err := d.Set("public_ip_address_ids", flattenNetworkSubResourceID(props.PublicIPAddresses)); err != nil {
			return fmt.Errorf("setting `public_ip_address_ids`: %+v", err)
		}

		if err := d.Set("public_ip_prefix_ids", flattenNetworkSubResourceID(props.PublicIPPrefixes)); err != nil {
			return fmt.Errorf("setting `public_ip_prefix_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
