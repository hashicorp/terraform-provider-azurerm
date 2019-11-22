package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznetwork "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubCreateUpdate,
		Read:   resourceArmVirtualHubRead,
		Update: resourceArmVirtualHubCreateUpdate,
		Delete: resourceArmVirtualHubDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetwork.ValidateVirtualHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.CIDR,
			},

			"virtual_wan_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"virtual_network_connection": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[\da-zA-Z][-_.\da-zA-Z]{0,78}[_\da-zA-Z]$`),
								`The name must be between 1 and 80 characters and begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.`,
							),
						},
						"remote_virtual_network_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"allow_hub_to_remote_vnet_transit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_remote_vnet_to_use_hub_vnet_gateways": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_internet_security": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"route": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_prefixes": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},
						"next_hop_ip_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
					},
				},
			},

			"s2s_vpn_gateway_id": {
				Type:         schema.TypeString,
				Computed:     true,
			},

			"p2s_vpn_gateway_id": {
				Type:         schema.TypeString,
				Computed:     true,
			},

			"express_route_gateway_id": {
				Type:         schema.TypeString,
				Computed:     true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	addressPrefix := d.Get("address_prefix").(string)
	virtualWanId := d.Get("virtual_wan_id").(string)
	virtualNetworkConnection := d.Get("virtual_network_connection").([]interface{})
	route := d.Get("route").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.VirtualHub{
		Location: utils.String(location),
		VirtualHubProperties: &network.VirtualHubProperties{
			AddressPrefix: utils.String(addressPrefix),
			VirtualWan: &network.SubResource{
				ID: &virtualWanId,
			},
			VirtualNetworkConnections: expandArmVirtualHubVirtualNetworkConnection(virtualNetworkConnection),
			RouteTable:                expandArmVirtualHubRoute(route),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualHubRead(d, meta)
}

func resourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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

func resourceArmVirtualHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmVirtualHubVirtualNetworkConnection(input []interface{}) *[]network.HubVirtualNetworkConnection {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.HubVirtualNetworkConnection, 0)

	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			name := v["name"].(string)
			remoteVirtualNetworkId := v["remote_virtual_network_id"].(string)
			allowHubToRemoteVnetTransit := v["allow_hub_to_remote_vnet_transit"].(bool)
			allowRemoteVnetToUseHubVnetGateways := v["allow_remote_vnet_to_use_hub_vnet_gateways"].(bool)
			enableInternetSecurity := v["enable_internet_security"].(bool)

			result := network.HubVirtualNetworkConnection{
				Name: utils.String(name),
				HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
					RemoteVirtualNetwork: &network.SubResource{
						ID: utils.String(remoteVirtualNetworkId),
					},
					AllowHubToRemoteVnetTransit:         utils.Bool(allowHubToRemoteVnetTransit),
					AllowRemoteVnetToUseHubVnetGateways: utils.Bool(allowRemoteVnetToUseHubVnetGateways),
					EnableInternetSecurity:              utils.Bool(enableInternetSecurity),
				},
			}

			results = append(results, result)
		}
	}

	return &results
}

func expandArmVirtualHubRoute(input []interface{}) *network.VirtualHubRouteTable {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.VirtualHubRoute, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			addressPrefixes := v["address_prefixes"].([]interface{})
			nextHopIpAddress := v["next_hop_ip_address"].(string)

			result := network.VirtualHubRoute{
				AddressPrefixes:  utils.ExpandStringSlice(addressPrefixes),
				NextHopIPAddress: utils.String(nextHopIpAddress),
			}

			results = append(results, result)
		}
	}

	result := network.VirtualHubRouteTable{
		Routes: &results,
	}

	return &result
}

func flattenArmVirtualHubVirtualNetworkConnection(input *[]network.HubVirtualNetworkConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		c := make(map[string]interface{})

		if name := item.Name; name != nil {
			c["name"] = *name
		}
		if props := item.HubVirtualNetworkConnectionProperties; props != nil {
			if v := props.RemoteVirtualNetwork; v != nil {
				c["remote_virtual_network_id"] = *v.ID
			}
			if v := props.AllowHubToRemoteVnetTransit; v != nil {
				c["allow_hub_to_remote_vnet_transit"] = *v
			}
			if v := props.AllowRemoteVnetToUseHubVnetGateways; v != nil {
				c["allow_remote_vnet_to_use_hub_vnet_gateways"] = *v
			}
			if v := props.EnableInternetSecurity; v != nil {
				c["enable_internet_security"] = *v
			}
		}

		results = append(results, c)
	}

	return results
}

func flattenArmVirtualHubRoute(input *[]network.VirtualHubRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		c := make(map[string]interface{})

		c["address_prefixes"] = utils.FlattenStringSlice(item.AddressPrefixes)
		if v := item.NextHopIPAddress; v != nil {
			c["next_hop_ip_address"] = *v
		}

		results = append(results, c)
	}

	return results
}
