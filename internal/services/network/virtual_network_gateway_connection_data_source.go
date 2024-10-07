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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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

			// TODO 4.0: change this from enable_* to *_enabled
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

			"private_link_fast_path_enabled": {
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
	client := meta.(*clients.Client).Network.VirtualNetworkGatewayConnections
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
		d.Set("enable_bgp", props.EnableBgp)
		d.Set("ingress_bytes_transferred", props.IngressBytesTransferred)
		d.Set("egress_bytes_transferred", props.EgressBytesTransferred)
		d.Set("use_policy_based_traffic_selectors", props.UsePolicyBasedTrafficSelectors)
		d.Set("express_route_gateway_bypass", props.ExpressRouteGatewayBypass)
		d.Set("private_link_fast_path_enabled", props.EnablePrivateLinkFastPath)
		d.Set("type", string(props.ConnectionType))
		d.Set("connection_protocol", string(pointer.From(props.ConnectionProtocol)))
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

		ipsecPoliciesSettingsFlat := flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(props.IPsecPolicies)
		if err := d.Set("ipsec_policy", ipsecPoliciesSettingsFlat); err != nil {
			return fmt.Errorf("setting `ipsec_policy`: %+v", err)
		}

		trafficSelectorsPolicyFlat := flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(props.TrafficSelectorPolicies)
		if err := d.Set("traffic_selector_policy", trafficSelectorsPolicyFlat); err != nil {
			return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
		}
	}

	return nil
}

func flattenVirtualNetworkGatewayConnectionDataSourceIpsecPolicies(ipsecPolicies *[]virtualnetworkgatewayconnections.IPsecPolicy) []interface{} {
	schemaIpsecPolicies := make([]interface{}, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]interface{})

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

func flattenVirtualNetworkGatewayConnectionDataSourcePolicyTrafficSelectors(trafficSelectorPolicies *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy) []interface{} {
	schemaTrafficSelectorPolicies := make([]interface{}, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]interface{}{
				"local_address_cidrs":  trafficSelectorPolicy.LocalAddressRanges,
				"remote_address_cidrs": trafficSelectorPolicy.RemoteAddressRanges,
			})
		}
	}

	return schemaTrafficSelectorPolicies
}
