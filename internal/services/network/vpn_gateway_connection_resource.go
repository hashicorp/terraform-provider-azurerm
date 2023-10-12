// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVPNGatewayConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVpnGatewayConnectionResourceCreateUpdate,
		Read:   resourceVpnGatewayConnectionResourceRead,
		Update: resourceVpnGatewayConnectionResourceCreateUpdate,
		Delete: resourceVpnGatewayConnectionResourceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseVPNConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vpn_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVpnGatewayID,
			},

			"remote_vpn_site_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVpnSiteID,
			},

			"internet_security_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Service will create a route table for the user if this is not specified.
			"routing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"associated_route_table": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.HubRouteTableID,
						},

						"inbound_route_map_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.RouteMapID,
						},

						"outbound_route_map_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.RouteMapID,
						},

						"propagated_route_table": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"route_table_ids": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validate.HubRouteTableID,
										},
									},

									"labels": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},
			},

			"vpn_link": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"vpn_site_link_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
							// The vpn site link associated with one link connection can not be updated
							ForceNew:     true,
							ValidateFunc: virtualwans.ValidateVpnSiteLinkID,
						},

						"egress_nat_rule_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.VpnGatewayNatRuleID,
							},
						},

						"ingress_nat_rule_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.VpnGatewayNatRuleID,
							},
						},

						"connection_mode": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnLinkConnectionMode(), false),
							Default:      string(virtualwans.VpnLinkConnectionModeDefault),
						},

						"route_weight": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},

						"protocol": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVirtualNetworkGatewayConnectionProtocol(), false),
							Default:      string(virtualwans.VirtualNetworkGatewayConnectionProtocolIKEvTwo),
						},

						"bandwidth_mbps": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      10,
						},

						"shared_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"bgp_enabled": {
							Type:     pluginsdk.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},

						"ipsec_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"sa_lifetime_sec": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(300, 172799),
									},
									"sa_data_size_kb": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1024, 2147483647),
									},
									"encryption_algorithm": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecEncryption(), false),
									},
									"integrity_algorithm": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecIntegrity(), false),
									},

									"ike_encryption_algorithm": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeEncryption(), false),
									},

									"ike_integrity_algorithm": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeIntegrity(), false),
									},

									"dh_group": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForDhGroup(), false),
									},

									"pfs_group": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForPfsGroup(), false),
									},
								},
							},
						},

						"ratelimit_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"local_azure_ip_address_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"policy_based_traffic_selector_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"custom_bgp_address": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"ip_address": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPv4Address,
									},

									"ip_configuration_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"traffic_selector_policy": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_address_ranges": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},

						"remote_address_ranges": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},
					},
				},
			},
		},
	}
}

func resourceVpnGatewayConnectionResourceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	gatewayId, err := virtualwans.ParseVpnGatewayID(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewVPNConnectionID(gatewayId.SubscriptionId, gatewayId.ResourceGroupName, gatewayId.VpnGatewayName, name)
	if d.IsNewResource() {
		resp, err := client.VpnConnectionsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_gateway_connection", id.ID())
		}
	}

	locks.ByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(gatewayId.VpnGatewayName, VPNGatewayResourceName)

	payload := virtualwans.VpnConnection{
		Properties: &virtualwans.VpnConnectionProperties{
			EnableInternetSecurity: pointer.To(d.Get("internet_security_enabled").(bool)),
			RemoteVpnSite: &virtualwans.SubResource{
				Id: pointer.To(d.Get("remote_vpn_site_id").(string)),
			},
			VpnLinkConnections:   expandVpnGatewayConnectionVpnSiteLinkConnections(d.Get("vpn_link").([]interface{})),
			RoutingConfiguration: expandVpnGatewayConnectionRoutingConfiguration(d.Get("routing").([]interface{})),
		},
	}

	if v, ok := d.GetOk("traffic_selector_policy"); ok {
		payload.Properties.TrafficSelectorPolicies = expandVpnGatewayConnectionTrafficSelectorPolicy(v.(*pluginsdk.Set).List())
	}

	if err := client.VpnConnectionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnGatewayConnectionResourceRead(d, meta)
}

func resourceVpnGatewayConnectionResourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVPNConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("vpn_gateway_id", virtualwans.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroupName, id.GatewayName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			vpnSiteId := ""
			if props.RemoteVpnSite != nil && props.RemoteVpnSite.Id != nil {
				theVpnSiteId, err := virtualwans.ParseVpnSiteIDInsensitively(*props.RemoteVpnSite.Id)
				if err != nil {
					return err
				}
				vpnSiteId = theVpnSiteId.ID()
			}
			d.Set("remote_vpn_site_id", vpnSiteId)

			enableInternetSecurity := false
			if props.EnableInternetSecurity != nil {
				enableInternetSecurity = *props.EnableInternetSecurity
			}
			d.Set("internet_security_enabled", enableInternetSecurity)

			if err := d.Set("routing", flattenVpnGatewayConnectionRoutingConfiguration(props.RoutingConfiguration)); err != nil {
				return fmt.Errorf(`setting "routing": %v`, err)
			}

			if err := d.Set("vpn_link", flattenVpnGatewayConnectionVpnSiteLinkConnections(props.VpnLinkConnections)); err != nil {
				return fmt.Errorf(`setting "vpn_link": %v`, err)
			}

			if err := d.Set("traffic_selector_policy", flattenVpnGatewayConnectionTrafficSelectorPolicy(props.TrafficSelectorPolicies)); err != nil {
				return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVpnGatewayConnectionResourceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVPNConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.GatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(id.GatewayName, VPNGatewayResourceName)

	if err := client.VpnConnectionsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVpnGatewayConnectionVpnSiteLinkConnections(input []interface{}) *[]virtualwans.VpnSiteLinkConnection {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.VpnSiteLinkConnection, 0)
	for _, itemRaw := range input {
		item := itemRaw.(map[string]interface{})
		v := virtualwans.VpnSiteLinkConnection{
			Name: utils.String(item["name"].(string)),
			Properties: &virtualwans.VpnSiteLinkConnectionProperties{
				VpnSiteLink: &virtualwans.SubResource{
					Id: utils.String(item["vpn_site_link_id"].(string)),
				},
				RoutingWeight:                  pointer.To(int64(item["route_weight"].(int))),
				VpnConnectionProtocolType:      pointer.To(virtualwans.VirtualNetworkGatewayConnectionProtocol(item["protocol"].(string))),
				VpnLinkConnectionMode:          pointer.To(virtualwans.VpnLinkConnectionMode(item["connection_mode"].(string))),
				ConnectionBandwidth:            pointer.To(int64(item["bandwidth_mbps"].(int))),
				EnableBgp:                      pointer.To(item["bgp_enabled"].(bool)),
				IPsecPolicies:                  expandVpnGatewayConnectionIpSecPolicies(item["ipsec_policy"].([]interface{})),
				EnableRateLimiting:             pointer.To(item["ratelimit_enabled"].(bool)),
				UseLocalAzureIPAddress:         pointer.To(item["local_azure_ip_address_enabled"].(bool)),
				UsePolicyBasedTrafficSelectors: pointer.To(item["policy_based_traffic_selector_enabled"].(bool)),
				VpnGatewayCustomBgpAddresses:   expandVpnGatewayConnectionCustomBgpAddresses(item["custom_bgp_address"].(*pluginsdk.Set).List()),
			},
		}

		if egressNatRuleIds := item["egress_nat_rule_ids"].(*pluginsdk.Set).List(); len(egressNatRuleIds) != 0 {
			v.Properties.EgressNatRules = expandVpnGatewayConnectionNatRuleIds(egressNatRuleIds)
		}

		if ingressNatRuleIds := item["ingress_nat_rule_ids"].(*pluginsdk.Set).List(); len(ingressNatRuleIds) != 0 {
			v.Properties.IngressNatRules = expandVpnGatewayConnectionNatRuleIds(ingressNatRuleIds)
		}

		if sharedKey := item["shared_key"]; sharedKey != "" {
			v.Properties.SharedKey = pointer.To(sharedKey.(string))
		}
		result = append(result, v)
	}

	return &result
}

func flattenVpnGatewayConnectionVpnSiteLinkConnections(input *[]virtualwans.VpnSiteLinkConnection) interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, item := range *input {
		if item.Properties == nil {
			continue
		}

		props := *item.Properties

		connectionProtocolType := ""
		if props.VpnConnectionProtocolType != nil {
			connectionProtocolType = string(*props.VpnConnectionProtocolType)
		}

		vpnLinkConnectionMode := ""
		if props.VpnLinkConnectionMode != nil {
			vpnLinkConnectionMode = string(*props.VpnLinkConnectionMode)
		}

		vpnSiteLinkId := ""
		if props.VpnSiteLink != nil && props.VpnSiteLink.Id != nil {
			vpnSiteLinkId = *props.VpnSiteLink.Id
		}

		output = append(output, map[string]interface{}{
			"name":                                  pointer.From(item.Name),
			"egress_nat_rule_ids":                   flattenVpnGatewayConnectionNatRuleIds(props.EgressNatRules),
			"ingress_nat_rule_ids":                  flattenVpnGatewayConnectionNatRuleIds(props.IngressNatRules),
			"vpn_site_link_id":                      vpnSiteLinkId,
			"route_weight":                          int(pointer.From(props.RoutingWeight)),
			"protocol":                              connectionProtocolType,
			"connection_mode":                       vpnLinkConnectionMode,
			"bandwidth_mbps":                        int(pointer.From(props.ConnectionBandwidth)),
			"shared_key":                            pointer.From(props.SharedKey),
			"bgp_enabled":                           pointer.From(props.EnableBgp),
			"ipsec_policy":                          flattenVpnGatewayConnectionIpSecPolicies(props.IPsecPolicies),
			"ratelimit_enabled":                     pointer.From(props.EnableRateLimiting),
			"local_azure_ip_address_enabled":        pointer.From(props.UseLocalAzureIPAddress),
			"policy_based_traffic_selector_enabled": pointer.From(props.UsePolicyBasedTrafficSelectors),
			"custom_bgp_address":                    flattenVpnGatewayConnectionCustomBgpAddresses(props.VpnGatewayCustomBgpAddresses),
		})
	}

	return output
}

func expandVpnGatewayConnectionIpSecPolicies(input []interface{}) *[]virtualwans.IPsecPolicy {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.IPsecPolicy, 0)
	for _, itemRaw := range input {
		item := itemRaw.(map[string]interface{})
		result = append(result, virtualwans.IPsecPolicy{
			SaLifeTimeSeconds:   int64(item["sa_lifetime_sec"].(int)),
			SaDataSizeKilobytes: int64(item["sa_data_size_kb"].(int)),
			IPsecEncryption:     virtualwans.IPsecEncryption(item["encryption_algorithm"].(string)),
			IPsecIntegrity:      virtualwans.IPsecIntegrity(item["integrity_algorithm"].(string)),
			IkeEncryption:       virtualwans.IkeEncryption(item["ike_encryption_algorithm"].(string)),
			IkeIntegrity:        virtualwans.IkeIntegrity(item["ike_integrity_algorithm"].(string)),
			DhGroup:             virtualwans.DhGroup(item["dh_group"].(string)),
			PfsGroup:            virtualwans.PfsGroup(item["pfs_group"].(string)),
		})
	}

	return &result
}

func flattenVpnGatewayConnectionIpSecPolicies(input *[]virtualwans.IPsecPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, item := range *input {
		output = append(output, map[string]interface{}{
			"sa_lifetime_sec":          int(item.SaLifeTimeSeconds),
			"sa_data_size_kb":          int(item.SaDataSizeKilobytes),
			"encryption_algorithm":     string(item.IPsecEncryption),
			"integrity_algorithm":      string(item.IPsecIntegrity),
			"ike_encryption_algorithm": string(item.IkeEncryption),
			"ike_integrity_algorithm":  string(item.IkeIntegrity),
			"dh_group":                 string(item.DhGroup),
			"pfs_group":                string(item.PfsGroup),
		})
	}

	return output
}

func expandVpnGatewayConnectionRoutingConfiguration(input []interface{}) *virtualwans.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &virtualwans.RoutingConfiguration{
		AssociatedRouteTable: &virtualwans.SubResource{
			Id: utils.String(raw["associated_route_table"].(string)),
		},
	}

	if inboundRouteMapId := raw["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		output.InboundRouteMap = &virtualwans.SubResource{
			Id: utils.String(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := raw["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		output.OutboundRouteMap = &virtualwans.SubResource{
			Id: utils.String(outboundRouteMapId),
		}
	}

	if v := raw["propagated_route_table"].([]interface{}); len(v) != 0 {
		output.PropagatedRouteTables = expandVpnGatewayConnectionPropagatedRouteTable(v)
	}

	return output
}

func flattenVpnGatewayConnectionRoutingConfiguration(input *virtualwans.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	associateRouteTable := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associateRouteTable = *input.AssociatedRouteTable.Id
	}

	var inboundRouteMapId string
	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		inboundRouteMapId = *input.InboundRouteMap.Id
	}

	var outboundRouteMapId string
	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		outboundRouteMapId = *input.OutboundRouteMap.Id
	}

	return []interface{}{
		map[string]interface{}{
			"propagated_route_table": flattenVpnGatewayConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"associated_route_table": associateRouteTable,
			"inbound_route_map_id":   inboundRouteMapId,
			"outbound_route_map_id":  outboundRouteMapId,
		},
	}
}

func flattenVpnGatewayConnectionPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	labels := make([]interface{}, 0)
	if input.Labels != nil {
		labels = utils.FlattenStringSlice(input.Labels)
	}

	routeTableIds := make([]interface{}, 0)
	if input.Ids != nil {
		for _, id := range *input.Ids {
			if id.Id == nil {
				continue
			}
			routeTableIds = append(routeTableIds, *id.Id)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func expandVpnGatewayConnectionTrafficSelectorPolicy(input []interface{}) *[]virtualwans.TrafficSelectorPolicy {
	results := make([]virtualwans.TrafficSelectorPolicy, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, virtualwans.TrafficSelectorPolicy{
			LocalAddressRanges:  pointer.From(utils.ExpandStringSlice(v["local_address_ranges"].(*pluginsdk.Set).List())),
			RemoteAddressRanges: pointer.From(utils.ExpandStringSlice(v["remote_address_ranges"].(*pluginsdk.Set).List())),
		})
	}

	return &results
}

func flattenVpnGatewayConnectionTrafficSelectorPolicy(input *[]virtualwans.TrafficSelectorPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"local_address_ranges":  utils.FlattenStringSlice(&item.LocalAddressRanges),
			"remote_address_ranges": utils.FlattenStringSlice(&item.RemoteAddressRanges),
		})
	}

	return results
}

func expandVpnGatewayConnectionPropagatedRouteTable(input []interface{}) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return &virtualwans.PropagatedRouteTable{}
	}

	v := input[0].(map[string]interface{})

	routeTableIds := make([]virtualwans.SubResource, 0)
	for _, val := range v["route_table_ids"].([]interface{}) {
		routeTableIds = append(routeTableIds, virtualwans.SubResource{
			Id: pointer.To(val.(string)),
		})
	}

	result := virtualwans.PropagatedRouteTable{
		Ids: pointer.To(routeTableIds),
	}
	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = utils.ExpandStringSlice(labels)
	}
	return &result
}

func expandVpnGatewayConnectionNatRuleIds(input []interface{}) *[]virtualwans.SubResource {
	results := make([]virtualwans.SubResource, 0)

	for _, item := range input {
		results = append(results, virtualwans.SubResource{
			Id: utils.String(item.(string)),
		})
	}

	return &results
}

func flattenVpnGatewayConnectionNatRuleIds(input *[]virtualwans.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var id string
		if item.Id != nil {
			id = *item.Id
		}

		results = append(results, id)
	}

	return results
}

func expandVpnGatewayConnectionCustomBgpAddresses(input []interface{}) *[]virtualwans.GatewayCustomBgpIPAddressIPConfiguration {
	results := make([]virtualwans.GatewayCustomBgpIPAddressIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, virtualwans.GatewayCustomBgpIPAddressIPConfiguration{
			CustomBgpIPAddress: v["ip_address"].(string),
			IPConfigurationId:  v["ip_configuration_id"].(string),
		})
	}

	return &results
}

func flattenVpnGatewayConnectionCustomBgpAddresses(input *[]virtualwans.GatewayCustomBgpIPAddressIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"ip_address":          item.CustomBgpIPAddress,
			"ip_configuration_id": item.IPConfigurationId,
		})
	}

	return results
}
