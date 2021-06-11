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

func dataSourceVirtualNetworkGatewayConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualNetworkGatewayConnectionRead,

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

			"location": azure.SchemaLocationForDataSource(),

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"shared_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"authorization_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dpd_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"enable_bgp": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"ingress_bytes_transferred": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"egress_bytes_transferred": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"use_policy_based_traffic_selectors": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"connection_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"routing_weight": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"virtual_network_gateway_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"peer_virtual_network_gateway_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"local_azure_ip_address_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"local_network_gateway_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"express_route_circuit_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"express_route_gateway_bypass": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"traffic_selector_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"remote_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"ipsec_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"sa_lifetime": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"sa_datasize": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"ipsec_encryption": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ipsec_integrity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ike_encryption": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ike_integrity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"dh_group": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"pfs_group": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Virtual Network Gateway Connection %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if resp.VirtualNetworkGatewayConnectionPropertiesFormat != nil {
		gwc := *resp.VirtualNetworkGatewayConnectionPropertiesFormat

		d.Set("shared_key", gwc.SharedKey)
		d.Set("authorization_key", gwc.AuthorizationKey)
		d.Set("enable_bgp", gwc.EnableBgp)
		d.Set("ingress_bytes_transferred", gwc.IngressBytesTransferred)
		d.Set("egress_bytes_transferred", gwc.EgressBytesTransferred)
		d.Set("use_policy_based_traffic_selectors", gwc.UsePolicyBasedTrafficSelectors)
		d.Set("express_route_gateway_bypass", gwc.ExpressRouteGatewayBypass)
		d.Set("type", string(gwc.ConnectionType))
		d.Set("connection_protocol", string(gwc.ConnectionProtocol))
		d.Set("routing_weight", gwc.RoutingWeight)

		if gwc.VirtualNetworkGateway1 != nil {
			d.Set("virtual_network_gateway_id", gwc.VirtualNetworkGateway1.ID)
		}

		if gwc.VirtualNetworkGateway2 != nil {
			d.Set("peer_virtual_network_gateway_id", gwc.VirtualNetworkGateway2.ID)
		}

		if gwc.LocalNetworkGateway2 != nil {
			d.Set("local_network_gateway_id", gwc.LocalNetworkGateway2.ID)
		}

		if gwc.Peer != nil {
			d.Set("express_route_circuit_id", gwc.Peer.ID)
		}

		if gwc.DpdTimeoutSeconds != nil {
			d.Set("dpd_timeout_seconds", gwc.DpdTimeoutSeconds)
		}

		if gwc.UseLocalAzureIPAddress != nil {
			d.Set("local_azure_ip_address_enabled", gwc.UseLocalAzureIPAddress)
		}

		d.Set("resource_guid", gwc.ResourceGUID)

		ipsecPoliciesSettingsFlat := flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(gwc.IpsecPolicies)
		if err := d.Set("ipsec_policy", ipsecPoliciesSettingsFlat); err != nil {
			return fmt.Errorf("Error setting `ipsec_policy`: %+v", err)
		}

		trafficSelectorsPolicyFlat := flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(gwc.TrafficSelectorPolicies)
		if err := d.Set("traffic_selector_policy", trafficSelectorsPolicyFlat); err != nil {
			return fmt.Errorf("Error setting `traffic_selector_policy`: %+v", err)
		}
	}

	return nil
}

func flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(ipsecPolicies *[]network.IpsecPolicy) []interface{} {
	schemaIpsecPolicies := make([]interface{}, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]interface{})

			schemaIpsecPolicy["dh_group"] = string(ipsecPolicy.DhGroup)
			schemaIpsecPolicy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
			schemaIpsecPolicy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
			schemaIpsecPolicy["ipsec_encryption"] = string(ipsecPolicy.IpsecEncryption)
			schemaIpsecPolicy["ipsec_integrity"] = string(ipsecPolicy.IpsecIntegrity)
			schemaIpsecPolicy["pfs_group"] = string(ipsecPolicy.PfsGroup)

			if ipsecPolicy.SaDataSizeKilobytes != nil {
				schemaIpsecPolicy["sa_datasize"] = int(*ipsecPolicy.SaDataSizeKilobytes)
			}

			if ipsecPolicy.SaLifeTimeSeconds != nil {
				schemaIpsecPolicy["sa_lifetime"] = int(*ipsecPolicy.SaLifeTimeSeconds)
			}

			schemaIpsecPolicies = append(schemaIpsecPolicies, schemaIpsecPolicy)
		}
	}

	return schemaIpsecPolicies
}

func flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(trafficSelectorPolicies *[]network.TrafficSelectorPolicy) []interface{} {
	schemaTrafficSelectorPolicies := make([]interface{}, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]interface{}{
				"local_address_cidrs":  utils.FlattenStringSlice(trafficSelectorPolicy.LocalAddressRanges),
				"remote_address_cidrs": utils.FlattenStringSlice(trafficSelectorPolicy.RemoteAddressRanges),
			})
		}
	}

	return schemaTrafficSelectorPolicies
}
