// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourcePointToSiteVPNGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePointToSiteVPNGatewayCreateUpdate,
		Read:   resourcePointToSiteVPNGatewayRead,
		Update: resourcePointToSiteVPNGatewayCreateUpdate,
		Delete: resourcePointToSiteVPNGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PointToSiteVpnGatewayID(id)
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
				ValidateFunc: networkValidate.VirtualHubID,
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
				// Code="P2SVpnGatewayCanHaveOnlyOneP2SConnectionConfiguration"
				// Message="Currently, P2SVpnGateway [ID] can have only one P2SConnectionConfiguration. Specified number of P2SConnectionConfiguration are 2.
				MaxItems: 1,
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
										ValidateFunc: networkValidate.HubRouteTableID,
									},

									"inbound_route_map_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: networkValidate.RouteMapID,
									},

									"outbound_route_map_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: networkValidate.RouteMapID,
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
														ValidateFunc: networkValidate.HubRouteTableID,
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

			"tags": tags.Schema(),
		},
	}
}

func resourcePointToSiteVPNGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPointToSiteVpnGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.P2sVpnGatewayName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_point_to_site_vpn_gateway", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnit := d.Get("scale_unit").(int)
	virtualHubId := d.Get("virtual_hub_id").(string)
	vpnServerConfigurationId := d.Get("vpn_server_configuration_id").(string)
	t := d.Get("tags").(map[string]interface{})

	connectionConfigurationsRaw := d.Get("connection_configuration").([]interface{})
	connectionConfigurations := expandPointToSiteVPNGatewayConnectionConfiguration(connectionConfigurationsRaw)

	parameters := network.P2SVpnGateway{
		Location: utils.String(location),
		P2SVpnGatewayProperties: &network.P2SVpnGatewayProperties{
			IsRoutingPreferenceInternet: utils.Bool(d.Get("routing_preference_internet_enabled").(bool)),
			P2SConnectionConfigurations: connectionConfigurations,
			VpnServerConfiguration: &network.SubResource{
				ID: utils.String(vpnServerConfigurationId),
			},
			VirtualHub: &network.SubResource{
				ID: utils.String(virtualHubId),
			},
			VpnGatewayScaleUnit: utils.Int32(int32(scaleUnit)),
		},
		Tags: tags.Expand(t),
	}
	customDNSServers := utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
	if len(*customDNSServers) != 0 {
		parameters.P2SVpnGatewayProperties.CustomDNSServers = customDNSServers
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.P2sVpnGatewayName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePointToSiteVPNGatewayRead(d, meta)
}

func resourcePointToSiteVPNGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PointToSiteVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.P2sVpnGatewayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.P2sVpnGatewayName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.P2SVpnGatewayProperties; props != nil {
		d.Set("dns_servers", utils.FlattenStringSlice(props.CustomDNSServers))
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
		if props.VirtualHub != nil && props.VirtualHub.ID != nil {
			virtualHubId = *props.VirtualHub.ID
		}
		d.Set("virtual_hub_id", virtualHubId)

		vpnServerConfigurationId := ""
		if props.VpnServerConfiguration != nil && props.VpnServerConfiguration.ID != nil {
			vpnServerConfigurationId = *props.VpnServerConfiguration.ID
		}
		d.Set("vpn_server_configuration_id", vpnServerConfigurationId)

		routingPreferenceInternetEnabled := false
		if props.IsRoutingPreferenceInternet != nil {
			routingPreferenceInternetEnabled = *props.IsRoutingPreferenceInternet
		}
		d.Set("routing_preference_internet_enabled", routingPreferenceInternetEnabled)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePointToSiteVPNGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PointToSiteVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.P2sVpnGatewayName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandPointToSiteVPNGatewayConnectionConfiguration(input []interface{}) *[]network.P2SConnectionConfiguration {
	configurations := make([]network.P2SConnectionConfiguration, 0)

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

		configurations = append(configurations, network.P2SConnectionConfiguration{
			Name: utils.String(name),
			P2SConnectionConfigurationProperties: &network.P2SConnectionConfigurationProperties{
				VpnClientAddressPool: &network.AddressSpace{
					AddressPrefixes: &addressPrefixes,
				},
				RoutingConfiguration:   expandPointToSiteVPNGatewayConnectionRouteConfiguration(raw["route"].([]interface{})),
				EnableInternetSecurity: utils.Bool(raw["internet_security_enabled"].(bool)),
			},
		})
	}

	return &configurations
}

func expandPointToSiteVPNGatewayConnectionRouteConfiguration(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	routingConfiguration := &network.RoutingConfiguration{
		AssociatedRouteTable: &network.SubResource{
			ID: utils.String(v["associated_route_table_id"].(string)),
		},
		PropagatedRouteTables: expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(v["propagated_route_table"].([]interface{})),
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		routingConfiguration.InboundRouteMap = &network.SubResource{
			ID: utils.String(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		routingConfiguration.OutboundRouteMap = &network.SubResource{
			ID: utils.String(outboundRouteMapId),
		}
	}

	return routingConfiguration
}

func expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	idRaws := utils.ExpandStringSlice(v["ids"].([]interface{}))
	ids := make([]network.SubResource, len(*idRaws))
	for i, item := range *idRaws {
		ids[i] = network.SubResource{
			ID: utils.String(item),
		}
	}
	return &network.PropagatedRouteTable{
		Labels: utils.ExpandStringSlice(v["labels"].(*pluginsdk.Set).List()),
		Ids:    &ids,
	}
}

func flattenPointToSiteVPNGatewayConnectionConfiguration(input *[]network.P2SConnectionConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		addressPrefixes := make([]interface{}, 0)
		enableInternetSecurity := false
		if props := v.P2SConnectionConfigurationProperties; props != nil {
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
		}

		output = append(output, map[string]interface{}{
			"name": name,
			"vpn_client_address_pool": []interface{}{
				map[string]interface{}{
					"address_prefixes": addressPrefixes,
				},
			},
			"route":                     flattenPointToSiteVPNGatewayConnectionRouteConfiguration(v.RoutingConfiguration),
			"internet_security_enabled": enableInternetSecurity,
		})
	}

	return output
}

func flattenPointToSiteVPNGatewayConnectionRouteConfiguration(input *network.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var associatedRouteTableId string
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.ID != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.ID
	}

	var inboundRouteMapId string
	if input.InboundRouteMap != nil && input.InboundRouteMap.ID != nil {
		inboundRouteMapId = *input.InboundRouteMap.ID
	}

	var outboundRouteMapId string
	if input.OutboundRouteMap != nil && input.OutboundRouteMap.ID != nil {
		outboundRouteMapId = *input.OutboundRouteMap.ID
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

func flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	ids := make([]string, 0)
	if input.Ids != nil {
		for _, item := range *input.Ids {
			if item.ID != nil {
				ids = append(ids, *item.ID)
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
