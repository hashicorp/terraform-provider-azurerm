// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/routefilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitPeeringCreate,
		Read:   resourceExpressRouteCircuitPeeringRead,
		Update: resourceExpressRouteCircuitPeeringUpdate,
		Delete: resourceExpressRouteCircuitPeeringDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseExpressRouteCircuitPeeringID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"peering_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePrivatePeering),
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering),
					string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering),
				}, false),
			},

			"express_route_circuit_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"primary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				RequiredWith: []string{
					"secondary_peer_address_prefix",
				},
			},

			"secondary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				RequiredWith: []string{
					"primary_peer_address_prefix",
				},
			},

			"ipv4_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"vlan_id": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"shared_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 25),
			},

			"peer_asn": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"microsoft_peering_config": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"advertised_public_prefixes": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},

						"customer_asn": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  0,
						},

						"routing_registry_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "NONE",
						},

						"advertised_communities": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"ipv6": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"microsoft_peering": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"advertised_public_prefixes": {
										Type:     pluginsdk.TypeList,
										MinItems: 1,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.IsCIDR,
										},
									},

									"customer_asn": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
										Default:  0,
									},

									"routing_registry_name": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "NONE",
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"advertised_communities": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},

						"primary_peer_address_prefix": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"secondary_peer_address_prefix": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"route_filter_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: routefilters.ValidateRouteFilterID,
						},
					},
				},
			},

			"azure_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_azure_port": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_azure_port": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"route_filter_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: routefilters.ValidateRouteFilterID,
			},

			"gateway_manager_etag": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceExpressRouteCircuitPeeringCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for Express Route Peering create.")

	id := commonids.NewExpressRouteCircuitPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_circuit_name").(string), d.Get("peering_type").(string))

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_express_route_circuit_peering", id.ID())
	}

	primaryPeerAddressPrefix := d.Get("primary_peer_address_prefix").(string)
	secondaryPeerAddressPrefix := d.Get("secondary_peer_address_prefix").(string)

	routeFilterId := d.Get("route_filter_id").(string)

	circuitConnClient := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	circuitConnectionId := commonids.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroupName, id.CircuitName, id.PeeringName)
	connResp, err := circuitConnClient.List(ctx, circuitConnectionId)
	if err != nil && !response.WasNotFound(connResp.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", circuitConnectionId, err)
	}

	var connection *[]expressroutecircuitpeerings.ExpressRouteCircuitConnection
	if model := connResp.Model; model != nil {
		connection = convertCircuitConnectionsToPeeringConnections(model)
	}

	parameters := expressroutecircuitpeerings.ExpressRouteCircuitPeering{
		Properties: &expressroutecircuitpeerings.ExpressRouteCircuitPeeringPropertiesFormat{
			PeeringType:        pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringType(id.PeeringName)),
			SharedKey:          pointer.To(d.Get("shared_key").(string)),
			AzureASN:           pointer.To(int64(d.Get("azure_asn").(int))),
			PeerASN:            pointer.To(int64(d.Get("peer_asn").(int))),
			VlanId:             pointer.To(int64(d.Get("vlan_id").(int))),
			GatewayManagerEtag: pointer.To(d.Get("gateway_manager_etag").(string)),
			Connections:        connection,
		},
	}

	ipv4Enabled := d.Get("ipv4_enabled").(bool)
	if ipv4Enabled {
		parameters.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)
	} else {
		parameters.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateDisabled)
	}

	if !strings.EqualFold(primaryPeerAddressPrefix, "") {
		parameters.Properties.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}

	if !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		parameters.Properties.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	if strings.EqualFold(id.PeeringName, string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering)) {
		peerings := d.Get("microsoft_peering_config").([]interface{})
		if len(peerings) == 0 && primaryPeerAddressPrefix != "" {
			return fmt.Errorf("`microsoft_peering_config` must be specified when config for Ipv4 and `peering_type` is set to `MicrosoftPeering`")
		}

		if len(peerings) != 0 && (primaryPeerAddressPrefix == "" || secondaryPeerAddressPrefix == "") {
			return fmt.Errorf("`primary_peer_address_prefix, secondary_peer_address_prefix` must be specified when config for Ipv4")
		}

		peeringConfig := expandExpressRouteCircuitPeeringMicrosoftConfig(peerings)
		parameters.Properties.MicrosoftPeeringConfig = peeringConfig

		if routeFilterId != "" {
			parameters.Properties.RouteFilter = &expressroutecircuitpeerings.SubResource{
				Id: pointer.To(routeFilterId),
			}
		}
	} else if routeFilterId != "" {
		return fmt.Errorf("`routeFilterId` may only be specified when `peering_type` is set to `MicrosoftPeering`")
	}

	ipv6Peering := d.Get("ipv6").([]interface{})
	if len(ipv6Peering) != 0 && id.PeeringName == string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering) {
		return fmt.Errorf("`ipv6` may only be specified when `peering_type` is `MicrosoftPeering` or `AzurePrivatePeering`")
	}

	ipv6PeeringConfig, err := expandExpressRouteCircuitIpv6PeeringConfig(ipv6Peering)
	if err != nil {
		return err
	}
	parameters.Properties.IPv6PeeringConfig = ipv6PeeringConfig

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitPeeringRead(d, meta)
}

func resourceExpressRouteCircuitPeeringUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Express Route Peering update.")

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	circuitConnClient := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	circuitConnectionId := commonids.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroupName, id.CircuitName, id.PeeringName)
	connResp, err := circuitConnClient.List(ctx, circuitConnectionId)
	if err != nil && !response.WasNotFound(connResp.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", circuitConnectionId, err)
	}

	var connection *[]expressroutecircuitpeerings.ExpressRouteCircuitConnection
	if model := connResp.Model; model != nil {
		connection = convertCircuitConnectionsToPeeringConnections(model)
	}

	payload.Properties.Connections = connection

	if d.HasChange("shared_key") {
		payload.Properties.SharedKey = pointer.To(d.Get("shared_key").(string))
	}

	if d.HasChange("azure_asn") {
		payload.Properties.AzureASN = pointer.To(int64(d.Get("azure_asn").(int)))
	}

	if d.HasChange("peer_asn") {
		payload.Properties.PeerASN = pointer.To(int64(d.Get("peer_asn").(int)))
	}

	if d.HasChange("vlan_id") {
		payload.Properties.VlanId = pointer.To(int64(d.Get("vlan_id").(int)))
	}

	if d.HasChange("gateway_manager_etag") {
		payload.Properties.GatewayManagerEtag = pointer.To(d.Get("gateway_manager_etag").(string))
	}

	if d.HasChange("ipv4_enabled") {

		ipv4Enabled := d.Get("ipv4_enabled").(bool)
		if ipv4Enabled {
			payload.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)
		} else {
			payload.Properties.State = pointer.To(expressroutecircuitpeerings.ExpressRoutePeeringStateDisabled)
		}
	}

	primaryPeerAddressPrefix := d.Get("primary_peer_address_prefix").(string)
	secondaryPeerAddressPrefix := d.Get("secondary_peer_address_prefix").(string)

	routeFilterId := d.Get("route_filter_id").(string)

	if d.HasChange("primary_peer_address_prefix") && !strings.EqualFold(primaryPeerAddressPrefix, "") {
		payload.Properties.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}

	if d.HasChange("secondary_peer_address_prefix") && !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		payload.Properties.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	if d.HasChange("microsoft_peering_config") {
		if strings.EqualFold(id.PeeringName, string(expressroutecircuitpeerings.ExpressRoutePeeringTypeMicrosoftPeering)) {
			peerings := d.Get("microsoft_peering_config").([]interface{})
			if len(peerings) == 0 && primaryPeerAddressPrefix != "" {
				return fmt.Errorf("`microsoft_peering_config` must be specified when config for Ipv4 and `peering_type` is set to `MicrosoftPeering`")
			}

			if len(peerings) != 0 && (primaryPeerAddressPrefix == "" || secondaryPeerAddressPrefix == "") {
				return fmt.Errorf("`primary_peer_address_prefix, secondary_peer_address_prefix` must be specified when config for Ipv4")
			}

			peeringConfig := expandExpressRouteCircuitPeeringMicrosoftConfig(peerings)
			payload.Properties.MicrosoftPeeringConfig = peeringConfig

			if d.HasChange("route_filter_id") && routeFilterId != "" {
				payload.Properties.RouteFilter = &expressroutecircuitpeerings.SubResource{
					Id: pointer.To(routeFilterId),
				}
			}
		} else if routeFilterId != "" {
			return fmt.Errorf("`routeFilterId` may only be specified when `peering_type` is set to `MicrosoftPeering`")
		}
	}

	if d.HasChange("ipv6") {
		ipv6Peering := d.Get("ipv6").([]interface{})
		if len(ipv6Peering) != 0 && id.PeeringName == string(expressroutecircuitpeerings.ExpressRoutePeeringTypeAzurePublicPeering) {
			return fmt.Errorf("`ipv6` may only be specified when `peering_type` is `MicrosoftPeering` or `AzurePrivatePeering`")
		}

		ipv6PeeringConfig, err := expandExpressRouteCircuitIpv6PeeringConfig(ipv6Peering)
		if err != nil {
			return err
		}
		payload.Properties.IPv6PeeringConfig = ipv6PeeringConfig
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitPeeringRead(d, meta)
}

func resourceExpressRouteCircuitPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("peering_type", id.PeeringName)
	d.Set("express_route_circuit_name", id.CircuitName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("azure_asn", props.AzureASN)
			d.Set("peer_asn", props.PeerASN)
			d.Set("primary_azure_port", props.PrimaryAzurePort)
			d.Set("secondary_azure_port", props.SecondaryAzurePort)
			d.Set("primary_peer_address_prefix", props.PrimaryPeerAddressPrefix)
			d.Set("secondary_peer_address_prefix", props.SecondaryPeerAddressPrefix)
			d.Set("vlan_id", props.VlanId)
			d.Set("gateway_manager_etag", props.GatewayManagerEtag)
			d.Set("ipv4_enabled", pointer.From(props.State) == expressroutecircuitpeerings.ExpressRoutePeeringStateEnabled)

			routeFilterId := ""
			if props.RouteFilter != nil && props.RouteFilter.Id != nil {
				routeFilterId = *props.RouteFilter.Id
			}
			d.Set("route_filter_id", routeFilterId)

			config := flattenExpressRouteCircuitPeeringMicrosoftConfig(props.MicrosoftPeeringConfig)
			if err := d.Set("microsoft_peering_config", config); err != nil {
				return fmt.Errorf("setting `microsoft_peering_config`: %+v", err)
			}
			if err := d.Set("ipv6", flattenExpressRouteCircuitIpv6PeeringConfig(props.IPv6PeeringConfig)); err != nil {
				return fmt.Errorf("setting `ipv6`: %+v", err)
			}
		}
	}

	return nil
}

func resourceExpressRouteCircuitPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitPeerings
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.CircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.CircuitName, expressRouteCircuitResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}

func expandExpressRouteCircuitPeeringMicrosoftConfig(input []interface{}) *expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig {
	if len(input) == 0 {
		return nil
	}
	peering := input[0].(map[string]interface{})

	prefixes := make([]string, 0)
	for _, v := range peering["advertised_public_prefixes"].([]interface{}) {
		prefixes = append(prefixes, v.(string))
	}

	advertisedCommunities := make([]string, 0)
	for _, v := range peering["advertised_communities"].([]interface{}) {
		advertisedCommunities = append(advertisedCommunities, v.(string))
	}

	return &expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig{
		AdvertisedPublicPrefixes: pointer.To(prefixes),
		CustomerASN:              pointer.To(int64(peering["customer_asn"].(int))),
		RoutingRegistryName:      pointer.To(peering["routing_registry_name"].(string)),
		AdvertisedCommunities:    pointer.To(advertisedCommunities),
	}
}

func expandExpressRouteCircuitIpv6PeeringConfig(input []interface{}) (*expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	peeringConfig := expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig{
		MicrosoftPeeringConfig: expandExpressRouteCircuitPeeringMicrosoftConfig(v["microsoft_peering"].([]interface{})),
		State:                  pointer.To(expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateEnabled),
	}

	primaryPeerAddressPrefix := v["primary_peer_address_prefix"].(string)
	secondaryPeerAddressPrefix := v["secondary_peer_address_prefix"].(string)
	if !strings.EqualFold(primaryPeerAddressPrefix, "") {
		peeringConfig.PrimaryPeerAddressPrefix = pointer.To(primaryPeerAddressPrefix)
	}
	if !strings.EqualFold(secondaryPeerAddressPrefix, "") {
		peeringConfig.SecondaryPeerAddressPrefix = pointer.To(secondaryPeerAddressPrefix)
	}

	ipv6Enabled := v["enabled"].(bool)
	if !ipv6Enabled {
		peeringConfig.State = pointer.To(expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateDisabled)
	}

	routeFilterId := v["route_filter_id"].(string)
	if routeFilterId != "" {
		if _, err := routefilters.ParseRouteFilterID(routeFilterId); err != nil {
			return nil, err
		}
		peeringConfig.RouteFilter = &expressroutecircuitpeerings.SubResource{
			Id: pointer.To(routeFilterId),
		}
	}
	return &peeringConfig, nil
}

func flattenExpressRouteCircuitPeeringMicrosoftConfig(input *expressroutecircuitpeerings.ExpressRouteCircuitPeeringConfig) interface{} {
	if input == nil {
		return []interface{}{}
	}

	config := make(map[string]interface{})
	prefixes := make([]string, 0)
	if customerASN := input.CustomerASN; customerASN != nil {
		config["customer_asn"] = *customerASN
	}
	if routingRegistryName := input.RoutingRegistryName; routingRegistryName != nil {
		config["routing_registry_name"] = *routingRegistryName
	}
	if ps := input.AdvertisedPublicPrefixes; ps != nil {
		prefixes = *ps
	}

	if community := input.AdvertisedCommunities; community != nil {
		config["advertised_communities"] = *community
	}
	config["advertised_public_prefixes"] = prefixes

	return []interface{}{config}
}

func flattenExpressRouteCircuitIpv6PeeringConfig(input *expressroutecircuitpeerings.IPv6ExpressRouteCircuitPeeringConfig) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var primaryPeerAddressPrefix string
	if input.PrimaryPeerAddressPrefix != nil {
		primaryPeerAddressPrefix = *input.PrimaryPeerAddressPrefix
	}
	var secondaryPeerAddressPrefix string
	if input.SecondaryPeerAddressPrefix != nil {
		secondaryPeerAddressPrefix = *input.SecondaryPeerAddressPrefix
	}
	routeFilterId := ""
	if input.RouteFilter != nil && input.RouteFilter.Id != nil {
		routeFilterId = *input.RouteFilter.Id
	}
	return []interface{}{
		map[string]interface{}{
			"microsoft_peering":             flattenExpressRouteCircuitPeeringMicrosoftConfig(input.MicrosoftPeeringConfig),
			"primary_peer_address_prefix":   primaryPeerAddressPrefix,
			"secondary_peer_address_prefix": secondaryPeerAddressPrefix,
			"route_filter_id":               routeFilterId,
			"enabled":                       pointer.From(input.State) == expressroutecircuitpeerings.ExpressRouteCircuitPeeringStateEnabled,
		},
	}
}

func convertCircuitConnectionsToPeeringConnections(input *[]expressroutecircuitconnections.ExpressRouteCircuitConnection) *[]expressroutecircuitpeerings.ExpressRouteCircuitConnection {
	output := make([]expressroutecircuitpeerings.ExpressRouteCircuitConnection, 0)

	if input == nil || len(*input) == 0 {
		return &output
	}

	for _, i := range *input {
		o := expressroutecircuitpeerings.ExpressRouteCircuitConnection{
			Etag: i.Etag,
			Id:   i.Id,
			Name: i.Name,
			Type: i.Type,
		}

		if props := i.Properties; props != nil {
			o.Properties = &expressroutecircuitpeerings.ExpressRouteCircuitConnectionPropertiesFormat{
				AddressPrefix:           props.AddressPrefix,
				AuthorizationKey:        props.AuthorizationKey,
				CircuitConnectionStatus: (*expressroutecircuitpeerings.CircuitConnectionStatus)(props.CircuitConnectionStatus),
				ExpressRouteCircuitPeering: &expressroutecircuitpeerings.SubResource{
					Id: props.ExpressRouteCircuitPeering.Id,
				},
				IPv6CircuitConnectionConfig: &expressroutecircuitpeerings.IPv6CircuitConnectionConfig{
					AddressPrefix:           props.IPv6CircuitConnectionConfig.AddressPrefix,
					CircuitConnectionStatus: (*expressroutecircuitpeerings.CircuitConnectionStatus)(props.IPv6CircuitConnectionConfig.CircuitConnectionStatus),
				},
				PeerExpressRouteCircuitPeering: &expressroutecircuitpeerings.SubResource{
					Id: props.PeerExpressRouteCircuitPeering.Id,
				},
				ProvisioningState: (*expressroutecircuitpeerings.ProvisioningState)(props.ProvisioningState),
			}
		}
		output = append(output, o)
	}
	return &output
}
