// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualHubConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualHubConnectionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VirtualHubConnectionName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"virtual_hub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.VirtualHubName,
			},

			"virtual_hub_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"remote_virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"internet_security_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"routing": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"associated_route_table_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"inbound_route_map_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"outbound_route_map_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"propagated_route_table": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"labels": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"route_table_ids": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},

						"static_vnet_local_route_override_criteria": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"static_vnet_propagate_static_routes_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						// lintignore:XS003
						"static_vnet_route": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"address_prefixes": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"next_hop_ip_address": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualHubConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewHubVirtualNetworkConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_hub_name").(string), d.Get("name").(string))

	resp, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HubVirtualNetworkConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("virtual_hub_name", id.VirtualHubName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("internet_security_enabled", props.EnableInternetSecurity)
			remoteVirtualNetworkId := ""
			if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.Id != nil {
				remoteVirtualNetworkId = *props.RemoteVirtualNetwork.Id
			}
			d.Set("remote_virtual_network_id", remoteVirtualNetworkId)

			if err := d.Set("routing", flattenVirtualHubConnectionRouting(props.RoutingConfiguration)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
		}
	}

	return nil
}
