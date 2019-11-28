package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	aznetwork "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualHubRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

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
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

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

		var expressRouteGatewayId *string
		if props.ExpressRouteGateway != nil {
			expressRouteGatewayId = props.ExpressRouteGateway.ID
		}
		d.Set("express_route_gateway_id", expressRouteGatewayId)

		var p2sVpnGatewayId *string
		if props.P2SVpnGateway != nil {
			p2sVpnGatewayId = props.P2SVpnGateway.ID
		}
		d.Set("p2s_vpn_gateway_id", p2sVpnGatewayId)

		if err := d.Set("route", flattenArmVirtualHubRoute(props.RouteTable)); err != nil {
			return fmt.Errorf("Error setting `route`: %+v", err)
		}

		var vpnGatewayId *string
		if props.VpnGateway != nil {
			vpnGatewayId = props.VpnGateway.ID
		}
		d.Set("s2s_vpn_gateway_id", vpnGatewayId)

		var virtualWanId *string
		if props.VirtualWan != nil {
			virtualWanId = props.VirtualWan.ID
		}
		d.Set("virtual_wan_id", virtualWanId)

		if err := d.Set("virtual_network_connection", flattenArmVirtualHubVirtualNetworkConnection(props.VirtualNetworkConnections)); err != nil {
			return fmt.Errorf("Error setting `virtual_network_connection`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
