// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Virtual Network Gateway Connection.",
		Read:        dataSourceVirtualNetworkGatewayConnectionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Virtual Network Gateway Connection.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"type": {
				Description: "The type of connection.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"shared_key": {
				Description: "The shared IPSec key.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"authorization_key": {
				Description: "The authorization key associated with the Express Route Circuit. This field is present only if the type is an ExpressRoute connection.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"dpd_timeout_seconds": {
				Description: "The dead peer detection timeout of this connection in seconds.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"bgp_enabled": {
				Description: "If `true`, BGP (Border Gateway Protocol) is enabled for this connection.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
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
				Description: "If `true`, policy-based traffic selectors are enabled for this connection. Enabling policy-based traffic selectors requires an `ipsec_policy` block.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"connection_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"routing_weight": {
				Description: "The routing weight.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"virtual_network_gateway_id": {
				Description: "The ID of the Virtual Network Gateway in which the connection is created.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"peer_virtual_network_gateway_id": {
				Description: "The ID of the peer Virtual Network Gateway when a VNet-to-VNet connection (i.e. when `type` is `Vnet2Vnet`).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"local_azure_ip_address_enabled": {
				Description: "Use private local Azure IP for the connection.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"local_network_gateway_id": {
				Description: "The ID of the Local Network Gateway when a Site-to-Site connection (i.e. when `type` is `IPsec`).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"express_route_circuit_id": {
				Description: "The ID of the Express Route Circuit (i.e. when `type` is `ExpressRoute`).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"express_route_gateway_bypass": {
				Description: "If `true`, data packets will bypass ExpressRoute Gateway for data forwarding. This is only valid for ExpressRoute connections.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"private_link_fast_path_enabled": {
				Description: "If `true`, data packets will bypass the Express Route gateway when accessing private-links. This is only valid for ExpressRoute connections, on the conditions described in [the relevant section in the Azure documentation](https://learn.microsoft.com/azure/expressroute/expressroute-howto-linkvnet-arm#fastpath-virtual-network-peering-user-defined-routes-udrs-and-private-link-support-for-expressroute-direct-connections).",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
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
							Description: "List of local CIDRs.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"remote_address_cidrs": {
							Description: "List of remote CIDRs.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
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
							Description: "The IPSec SA lifetime in seconds.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"sa_datasize": {
							Description: "The IPSec SA payload size in KB.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"ipsec_encryption": {
							Description: "The IPSec encryption algorithm.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"ipsec_integrity": {
							Description: "The IPSec integrity algorithm.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"ike_encryption": {
							Description: "The IKE encryption algorithm.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"ike_integrity": {
							Description: "The IKE integrity algorithm.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"dh_group": {
							Description: "The DH group used in IKE phase 1 for initial SA.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"pfs_group": {
							Description: "The DH group used in IKE phase 2 for new child SA.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGatewayConnections
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkgatewayconnections.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)

	respKey, err := client.GetSharedKey(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving Shared Key for %s: %+v", id, err)
	}

	if model := respKey.Model; model != nil {
		d.Set("shared_key", model.Value)
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties

		d.Set("authorization_key", props.AuthorizationKey)
		d.Set("bgp_enabled", props.EnableBgp)
		d.Set("ingress_bytes_transferred", props.IngressBytesTransferred)
		d.Set("egress_bytes_transferred", props.EgressBytesTransferred)
		d.Set("use_policy_based_traffic_selectors", props.UsePolicyBasedTrafficSelectors)
		d.Set("express_route_gateway_bypass", props.ExpressRouteGatewayBypass)
		d.Set("private_link_fast_path_enabled", props.EnablePrivateLinkFastPath)
		d.Set("type", string(props.ConnectionType))
		d.Set("connection_protocol", pointer.FromEnum(props.ConnectionProtocol))
		d.Set("routing_weight", props.RoutingWeight)

		d.Set("virtual_network_gateway_id", props.VirtualNetworkGateway1.Id)

		if props.VirtualNetworkGateway2 != nil {
			d.Set("peer_virtual_network_gateway_id", props.VirtualNetworkGateway2.Id)
		}

		if props.LocalNetworkGateway2 != nil {
			d.Set("local_network_gateway_id", props.LocalNetworkGateway2.Id)
		}

		if props.Peer != nil {
			d.Set("express_route_circuit_id", props.Peer.Id)
		}

		if props.DpdTimeoutSeconds != nil {
			d.Set("dpd_timeout_seconds", props.DpdTimeoutSeconds)
		}

		if props.UseLocalAzureIPAddress != nil {
			d.Set("local_azure_ip_address_enabled", props.UseLocalAzureIPAddress)
		}

		d.Set("resource_guid", props.ResourceGuid)

		if err := d.Set("ipsec_policy", flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(props.IPsecPolicies)); err != nil {
			return fmt.Errorf("setting `ipsec_policy`: %+v", err)
		}

		if err := d.Set("traffic_selector_policy", flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(props.TrafficSelectorPolicies)); err != nil {
			return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
		}
	}

	return nil
}

func flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(ipsecPolicies *[]virtualnetworkgatewayconnections.IPsecPolicy) []any {
	schemaIpsecPolicies := make([]any, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]any)

			schemaIpsecPolicy["dh_group"] = string(ipsecPolicy.DhGroup)
			schemaIpsecPolicy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
			schemaIpsecPolicy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
			schemaIpsecPolicy["ipsec_encryption"] = string(ipsecPolicy.IPsecEncryption)
			schemaIpsecPolicy["ipsec_integrity"] = string(ipsecPolicy.IPsecIntegrity)
			schemaIpsecPolicy["pfs_group"] = string(ipsecPolicy.PfsGroup)
			schemaIpsecPolicy["sa_datasize"] = int(ipsecPolicy.SaDataSizeKilobytes)
			schemaIpsecPolicy["sa_lifetime"] = int(ipsecPolicy.SaLifeTimeSeconds)

			schemaIpsecPolicies = append(schemaIpsecPolicies, schemaIpsecPolicy)
		}
	}

	return schemaIpsecPolicies
}

func flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(trafficSelectorPolicies *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy) []any {
	schemaTrafficSelectorPolicies := make([]any, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]any{
				"local_address_cidrs":  trafficSelectorPolicy.LocalAddressRanges,
				"remote_address_cidrs": trafficSelectorPolicy.RemoteAddressRanges,
			})
		}
	}

	return schemaTrafficSelectorPolicies
}
