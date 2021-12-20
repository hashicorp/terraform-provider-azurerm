package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceExpressRouteCircuitPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitPeeringCreateUpdate,
		Read:   resourceExpressRouteCircuitPeeringRead,
		Update: resourceExpressRouteCircuitPeeringCreateUpdate,
		Delete: resourceExpressRouteCircuitPeeringDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRouteCircuitPeeringID(id)
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
					string(network.ExpressRoutePeeringTypeAzurePrivatePeering),
					string(network.ExpressRoutePeeringTypeAzurePublicPeering),
					string(network.ExpressRoutePeeringTypeMicrosoftPeering),
				}, false),
			},

			"express_route_circuit_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"primary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"secondary_peer_address_prefix": {
				Type:     pluginsdk.TypeString,
				Required: true,
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
							Required: true,
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

						"route_filter_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
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
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceExpressRouteCircuitPeeringCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for Express Route Peering create/update.")

	id := parse.NewExpressRouteCircuitPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_circuit_name").(string), d.Get("peering_type").(string))

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_peering", *existing.ID)
		}
	}

	sharedKey := d.Get("shared_key").(string)
	primaryPeerAddressPrefix := d.Get("primary_peer_address_prefix").(string)
	secondaryPeerAddressPrefix := d.Get("secondary_peer_address_prefix").(string)
	vlanId := d.Get("vlan_id").(int)
	azureASN := d.Get("azure_asn").(int)
	peerASN := d.Get("peer_asn").(int)
	route_filter_id := d.Get("route_filter_id").(string)

	parameters := network.ExpressRouteCircuitPeering{
		ExpressRouteCircuitPeeringPropertiesFormat: &network.ExpressRouteCircuitPeeringPropertiesFormat{
			PeeringType:                network.ExpressRoutePeeringType(id.PeeringName),
			SharedKey:                  utils.String(sharedKey),
			PrimaryPeerAddressPrefix:   utils.String(primaryPeerAddressPrefix),
			SecondaryPeerAddressPrefix: utils.String(secondaryPeerAddressPrefix),
			AzureASN:                   utils.Int32(int32(azureASN)),
			PeerASN:                    utils.Int64(int64(peerASN)),
			VlanID:                     utils.Int32(int32(vlanId)),
		},
	}

	if strings.EqualFold(id.PeeringName, string(network.ExpressRoutePeeringTypeMicrosoftPeering)) {
		peerings := d.Get("microsoft_peering_config").([]interface{})
		if len(peerings) == 0 {
			return fmt.Errorf("`microsoft_peering_config` must be specified when `peering_type` is set to `MicrosoftPeering`")
		}

		peeringConfig := expandExpressRouteCircuitPeeringMicrosoftConfig(peerings)
		parameters.ExpressRouteCircuitPeeringPropertiesFormat.MicrosoftPeeringConfig = peeringConfig

		if route_filter_id != "" {
			parameters.ExpressRouteCircuitPeeringPropertiesFormat.RouteFilter = &network.SubResource{
				ID: utils.String(route_filter_id),
			}
		}

		ipv6Peering := d.Get("ipv6").([]interface{})
		ipv6PeeringConfig, err := expandExpressRouteCircuitIpv6PeeringConfig(ipv6Peering)
		if err != nil {
			return err
		}
		parameters.ExpressRouteCircuitPeeringPropertiesFormat.Ipv6PeeringConfig = ipv6PeeringConfig
	} else {
		if route_filter_id != "" {
			return fmt.Errorf("`route_filter_id` may only be specified when `peering_type` is set to `MicrosoftPeering`")
		}

		ipv6Peering := d.Get("ipv6").([]interface{})
		if len(ipv6Peering) != 0 {
			return fmt.Errorf("`ipv6` may only be specified when `peering_type` is set to `MicrosoftPeering`")
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitPeeringRead(d, meta)
}

func resourceExpressRouteCircuitPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("peering_type", id.PeeringName)
	d.Set("express_route_circuit_name", id.ExpressRouteCircuitName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ExpressRouteCircuitPeeringPropertiesFormat; props != nil {
		d.Set("azure_asn", props.AzureASN)
		d.Set("peer_asn", props.PeerASN)
		d.Set("primary_azure_port", props.PrimaryAzurePort)
		d.Set("secondary_azure_port", props.SecondaryAzurePort)
		d.Set("primary_peer_address_prefix", props.PrimaryPeerAddressPrefix)
		d.Set("secondary_peer_address_prefix", props.SecondaryPeerAddressPrefix)
		d.Set("vlan_id", props.VlanID)

		routeFilterId := ""
		if props.RouteFilter != nil && props.RouteFilter.ID != nil {
			routeFilterId = *props.RouteFilter.ID
		}
		d.Set("route_filter_id", routeFilterId)

		config := flattenExpressRouteCircuitPeeringMicrosoftConfig(props.MicrosoftPeeringConfig)
		if err := d.Set("microsoft_peering_config", config); err != nil {
			return fmt.Errorf("setting `microsoft_peering_config`: %+v", err)
		}
		if err := d.Set("ipv6", flattenExpressRouteCircuitIpv6PeeringConfig(props.Ipv6PeeringConfig)); err != nil {
			return fmt.Errorf("setting `ipv6`: %+v", err)
		}
	}

	return nil
}

func resourceExpressRouteCircuitPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}

func expandExpressRouteCircuitPeeringMicrosoftConfig(input []interface{}) *network.ExpressRouteCircuitPeeringConfig {
	if len(input) == 0 {
		return nil
	}
	peering := input[0].(map[string]interface{})

	prefixes := make([]string, 0)
	inputPrefixes := peering["advertised_public_prefixes"].([]interface{})
	inputCustomerASN := int32(peering["customer_asn"].(int))
	inputRoutingRegistryName := peering["routing_registry_name"].(string)

	for _, v := range inputPrefixes {
		prefixes = append(prefixes, v.(string))
	}

	return &network.ExpressRouteCircuitPeeringConfig{
		AdvertisedPublicPrefixes: &prefixes,
		CustomerASN:              &inputCustomerASN,
		RoutingRegistryName:      &inputRoutingRegistryName,
	}
}

func expandExpressRouteCircuitIpv6PeeringConfig(input []interface{}) (*network.Ipv6ExpressRouteCircuitPeeringConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	peeringConfig := network.Ipv6ExpressRouteCircuitPeeringConfig{
		PrimaryPeerAddressPrefix:   utils.String(v["primary_peer_address_prefix"].(string)),
		SecondaryPeerAddressPrefix: utils.String(v["secondary_peer_address_prefix"].(string)),
		MicrosoftPeeringConfig:     expandExpressRouteCircuitPeeringMicrosoftConfig(v["microsoft_peering"].([]interface{})),
	}
	routeFilterId := v["route_filter_id"].(string)
	if routeFilterId != "" {
		if _, err := parse.RouteFilterID(routeFilterId); err != nil {
			return nil, err
		}
		peeringConfig.RouteFilter = &network.SubResource{
			ID: utils.String(routeFilterId),
		}
	}
	return &peeringConfig, nil
}

func flattenExpressRouteCircuitPeeringMicrosoftConfig(input *network.ExpressRouteCircuitPeeringConfig) interface{} {
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
	config["advertised_public_prefixes"] = prefixes

	return []interface{}{config}
}

func flattenExpressRouteCircuitIpv6PeeringConfig(input *network.Ipv6ExpressRouteCircuitPeeringConfig) []interface{} {
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
	if input.RouteFilter != nil && input.RouteFilter.ID != nil {
		routeFilterId = *input.RouteFilter.ID
	}
	return []interface{}{
		map[string]interface{}{
			"microsoft_peering":             flattenExpressRouteCircuitPeeringMicrosoftConfig(input.MicrosoftPeeringConfig),
			"primary_peer_address_prefix":   primaryPeerAddressPrefix,
			"secondary_peer_address_prefix": secondaryPeerAddressPrefix,
			"route_filter_id":               routeFilterId,
		},
	}
}
