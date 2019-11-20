package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	aznetwork "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualHubRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetwork.ValidateVirtualHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_wan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"s2s_vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"p2s_vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"express_route_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_network_connection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_virtual_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_hub_to_remote_vnet_transit": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allow_remote_vnet_to_use_hub_vnet_gateways": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_internet_security": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"next_hop_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Hub %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VirtualHubProperties; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		if props.VirtualWan != nil {
			if err := d.Set("virtual_wan_id", props.VirtualWan.ID); err != nil {
				return fmt.Errorf("Error setting `virtual_wan_id`: %+v", err)
			}
		}
		if props.VpnGateway != nil {
			if err := d.Set("s2s_vpn_gateway_id", props.VpnGateway.ID); err != nil {
				return fmt.Errorf("Error setting `s2s_vpn_gateway_id`: %+v", err)
			}
		}
		if props.P2SVpnGateway != nil {
			if err := d.Set("p2s_vpn_gateway_id", props.P2SVpnGateway.ID); err != nil {
				return fmt.Errorf("Error setting `p2s_vpn_gateway_id`: %+v", err)
			}
		}
		if props.ExpressRouteGateway != nil {
			if err := d.Set("express_route_gateway_id", props.ExpressRouteGateway.ID); err != nil {
				return fmt.Errorf("Error setting `express_route_gateway_id`: %+v", err)
			}
		}
		if err := d.Set("virtual_network_connection", flattenArmVirtualHubVirtualNetworkConnection(props.VirtualNetworkConnections)); err != nil {
			return fmt.Errorf("Error setting `virtual_network_connection`: %+v", err)
		}
		if props.RouteTable != nil {
			if err := d.Set("route", flattenArmVirtualHubRoute(props.RouteTable.Routes)); err != nil {
				return fmt.Errorf("Error setting `route`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
