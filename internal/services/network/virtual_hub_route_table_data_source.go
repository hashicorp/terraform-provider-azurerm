// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceVirtualHubRouteTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualHubRouteTableRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HubRouteTableID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: networkValidate.HubRouteTableName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"virtual_hub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: networkValidate.VirtualHubName,
			},

			"virtual_hub_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"labels": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"route": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"destinations": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"destinations_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"next_hop": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"next_hop_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualHubRouteTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewHubRouteTableID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_hub_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("virtual_hub_name", id.VirtualHubName)
	d.Set("virtual_hub_id", parse.NewVirtualHubID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName).ID())

	if props := resp.HubRouteTableProperties; props != nil {
		d.Set("labels", utils.FlattenStringSlice(props.Labels))

		if err := d.Set("route", flattenVirtualHubRouteTableHubRoutes(props.Routes)); err != nil {
			return fmt.Errorf("setting `route`: %+v", err)
		}
	}
	return nil
}
