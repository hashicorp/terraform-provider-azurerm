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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePointToSiteVPNGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePointToSiteVPNGatewayCreate,
		Read:   resourcePointToSiteVPNGatewayRead,
		Update: resourcePointToSiteVPNGatewayUpdate,
		Delete: resourcePointToSiteVPNGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseVirtualWANP2SVPNGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"vpn_server_configuration_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVpnServerConfigurationID,
			},

			"connection_configuration": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"vpn_client_address_pool": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"address_prefixes": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validate.CIDR,
										},
									},
								},
							},
						},

						"route": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"associated_route_table_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: virtualwans.ValidateHubRouteTableID,
									},

									"inbound_route_map_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: virtualwans.ValidateRouteMapID,
									},

									"outbound_route_map_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: virtualwans.ValidateRouteMapID,
									},

									"propagated_route_table": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"ids": {
													Type:     pluginsdk.TypeList,
													Required: true,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeString,
														ValidateFunc: virtualwans.ValidateHubRouteTableID,
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
						"internet_security_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
					},
				},
			},

			"scale_unit": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"routing_preference_internet_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"dns_servers": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePointToSiteVPNGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualWANP2SVPNGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.P2sVpnGatewaysGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_point_to_site_vpn_gateway", id.ID())
	}

	parameters := virtualwans.P2SVpnGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.P2SVpnGatewayProperties{
			IsRoutingPreferenceInternet: pointer.To(d.Get("routing_preference_internet_enabled").(bool)),
			P2SConnectionConfigurations: expandPointToSiteVPNGatewayConnectionConfiguration(d.Get("connection_configuration").([]interface{})),
			VpnServerConfiguration: &virtualwans.SubResource{
				Id: pointer.To(d.Get("vpn_server_configuration_id").(string)),
			},
			VirtualHub: &virtualwans.SubResource{
				Id: pointer.To(d.Get("virtual_hub_id").(string)),
			},
			VpnGatewayScaleUnit: pointer.To(int64(d.Get("scale_unit").(int))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	customDNSServers := utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
	if len(*customDNSServers) != 0 {
		parameters.Properties.CustomDnsServers = customDNSServers
	}

	if err := client.P2sVpnGatewaysCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePointToSiteVPNGatewayRead(d, meta)
}

func resourcePointToSiteVPNGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}
	existing, err := client.P2sVpnGatewaysGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	props := virtualwans.P2SVpnGatewayProperties{}
	if existing.Model.Properties != nil {
		props = *existing.Model.Properties
	}

	if d.HasChange("connection_configuration") {
		props.P2SConnectionConfigurations = expandPointToSiteVPNGatewayConnectionConfiguration(d.Get("connection_configuration").([]interface{}))
	}

	if d.HasChange("scale_unit") {
		props.VpnGatewayScaleUnit = pointer.To(int64(d.Get("scale_unit").(int)))
	}

	if d.HasChange("dns_servers") {
		customDNSServers := utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
		if len(*customDNSServers) != 0 {
			props.CustomDnsServers = customDNSServers
		}
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	existing.Model.Properties = &props

	if err := client.P2sVpnGatewaysCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePointToSiteVPNGatewayRead(d, meta)
}

func resourcePointToSiteVPNGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.P2sVpnGatewaysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.GatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {

			d.Set("dns_servers", utils.FlattenStringSlice(props.CustomDnsServers))
			flattenedConfigurations := flattenPointToSiteVPNGatewayConnectionConfiguration(props.P2SConnectionConfigurations)
			if err := d.Set("connection_configuration", flattenedConfigurations); err != nil {
				return fmt.Errorf("setting `connection_configuration`: %+v", err)
			}

			scaleUnit := 0
			if props.VpnGatewayScaleUnit != nil {
				scaleUnit = int(*props.VpnGatewayScaleUnit)
			}
			d.Set("scale_unit", scaleUnit)

			virtualHubId := ""
			if props.VirtualHub != nil && props.VirtualHub.Id != nil {
				virtualHubId = *props.VirtualHub.Id
			}
			d.Set("virtual_hub_id", virtualHubId)

			vpnServerConfigurationId := ""
			if props.VpnServerConfiguration != nil && props.VpnServerConfiguration.Id != nil {
				vpnServerConfigurationId = *props.VpnServerConfiguration.Id
			}
			d.Set("vpn_server_configuration_id", vpnServerConfigurationId)

			routingPreferenceInternetEnabled := false
			if props.IsRoutingPreferenceInternet != nil {
				routingPreferenceInternetEnabled = *props.IsRoutingPreferenceInternet
			}
			d.Set("routing_preference_internet_enabled", routingPreferenceInternetEnabled)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourcePointToSiteVPNGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.P2sVpnGatewaysDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandPointToSiteVPNGatewayConnectionConfiguration(input []interface{}) *[]virtualwans.P2SConnectionConfiguration {
	configurations := make([]virtualwans.P2SConnectionConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		addressPrefixes := make([]string, 0)
		name := raw["name"].(string)

		clientAddressPoolsRaw := raw["vpn_client_address_pool"].([]interface{})
		for _, clientV := range clientAddressPoolsRaw {
			clientRaw := clientV.(map[string]interface{})

			addressPrefixesRaw := clientRaw["address_prefixes"].(*pluginsdk.Set).List()
			for _, prefix := range addressPrefixesRaw {
				addressPrefixes = append(addressPrefixes, prefix.(string))
			}
		}

		configurations = append(configurations, virtualwans.P2SConnectionConfiguration{
			Name: pointer.To(name),
			Properties: &virtualwans.P2SConnectionConfigurationProperties{
				VpnClientAddressPool: &virtualwans.AddressSpace{
					AddressPrefixes: &addressPrefixes,
				},
				RoutingConfiguration:   expandPointToSiteVPNGatewayConnectionRouteConfiguration(raw["route"].([]interface{})),
				EnableInternetSecurity: pointer.To(raw["internet_security_enabled"].(bool)),
			},
		})
	}

	return &configurations
}

func expandPointToSiteVPNGatewayConnectionRouteConfiguration(input []interface{}) *virtualwans.RoutingConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	routingConfiguration := &virtualwans.RoutingConfiguration{
		AssociatedRouteTable: &virtualwans.SubResource{
			Id: pointer.To(v["associated_route_table_id"].(string)),
		},
		PropagatedRouteTables: expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(v["propagated_route_table"].([]interface{})),
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		routingConfiguration.InboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		routingConfiguration.OutboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	return routingConfiguration
}

func expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input []interface{}) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	idRaws := utils.ExpandStringSlice(v["ids"].([]interface{}))
	ids := make([]virtualwans.SubResource, len(*idRaws))
	for i, item := range *idRaws {
		ids[i] = virtualwans.SubResource{
			Id: pointer.To(item),
		}
	}
	return &virtualwans.PropagatedRouteTable{
		Labels: utils.ExpandStringSlice(v["labels"].(*pluginsdk.Set).List()),
		Ids:    &ids,
	}
}

func flattenPointToSiteVPNGatewayConnectionConfiguration(input *[]virtualwans.P2SConnectionConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		route := make([]interface{}, 0)
		addressPrefixes := make([]interface{}, 0)
		enableInternetSecurity := false
		if props := v.Properties; props != nil {
			if props.VpnClientAddressPool == nil {
				continue
			}

			if props.VpnClientAddressPool.AddressPrefixes != nil {
				for _, prefix := range *props.VpnClientAddressPool.AddressPrefixes {
					addressPrefixes = append(addressPrefixes, prefix)
				}
			}

			if props.EnableInternetSecurity != nil {
				enableInternetSecurity = *props.EnableInternetSecurity
			}

			if props.RoutingConfiguration != nil {
				route = flattenPointToSiteVPNGatewayConnectionRouteConfiguration(props.RoutingConfiguration)
			}
		}

		output = append(output, map[string]interface{}{
			"name": name,
			"vpn_client_address_pool": []interface{}{
				map[string]interface{}{
					"address_prefixes": addressPrefixes,
				},
			},
			"route":                     route,
			"internet_security_enabled": enableInternetSecurity,
		})
	}

	return output
}

func flattenPointToSiteVPNGatewayConnectionRouteConfiguration(input *virtualwans.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var associatedRouteTableId string
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
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
			"associated_route_table_id": associatedRouteTableId,
			"inbound_route_map_id":      inboundRouteMapId,
			"outbound_route_map_id":     outboundRouteMapId,
			"propagated_route_table":    flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input.PropagatedRouteTables),
		},
	}
}

func flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	ids := make([]string, 0)
	if input.Ids != nil {
		for _, item := range *input.Ids {
			if item.Id != nil {
				ids = append(ids, *item.Id)
			}
		}
	}
	return []interface{}{
		map[string]interface{}{
			"ids":    ids,
			"labels": utils.FlattenStringSlice(input.Labels),
		},
	}
}
