package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmExpressRouteCircuitPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRouteCircuitPeeringCreateUpdate,
		Read:   resourceArmExpressRouteCircuitPeeringRead,
		Update: resourceArmExpressRouteCircuitPeeringCreateUpdate,
		Delete: resourceArmExpressRouteCircuitPeeringDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"peering_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzurePrivatePeering),
					string(network.AzurePublicPeering),
					string(network.MicrosoftPeering),
				}, false),
			},

			"express_route_circuit_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"primary_peer_address_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"secondary_peer_address_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"shared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 25),
			},

			"peer_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"microsoft_peering_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_public_prefixes": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"customer_asn": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"routing_registry_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "NONE",
						},
					},
				},
			},

			"ipv6": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"microsoft_peering": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"advertised_public_prefixes": {
										Type:     schema.TypeList,
										MinItems: 1,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsCIDR,
										},
									},

									"customer_asn": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},

									"routing_registry_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "NONE",
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"primary_peer_address_prefix": {
							Type:     schema.TypeString,
							Required: true,
						},

						"secondary_peer_address_prefix": {
							Type:     schema.TypeString,
							Required: true,
						},

						"route_filter_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"azure_asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_azure_port": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_azure_port": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"route_filter_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmExpressRouteCircuitPeeringCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Express Route Peering create/update.")

	peeringType := d.Get("peering_type").(string)
	circuitName := d.Get("express_route_circuit_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, circuitName, peeringType)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Peering %q (ExpressRoute Circuit %q / Resource Group %q): %s", peeringType, circuitName, resourceGroup, err)
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
			PeeringType:                network.ExpressRoutePeeringType(peeringType),
			SharedKey:                  utils.String(sharedKey),
			PrimaryPeerAddressPrefix:   utils.String(primaryPeerAddressPrefix),
			SecondaryPeerAddressPrefix: utils.String(secondaryPeerAddressPrefix),
			AzureASN:                   utils.Int32(int32(azureASN)),
			PeerASN:                    utils.Int64(int64(peerASN)),
			VlanID:                     utils.Int32(int32(vlanId)),
		},
	}

	if strings.EqualFold(peeringType, string(network.MicrosoftPeering)) {
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
	} else if route_filter_id != "" {
		return fmt.Errorf("`route_filter_id` may only be specified when `peering_type` is set to `MicrosoftPeering`")
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, circuitName, peeringType, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, circuitName, peeringType)
	if err != nil {
		return err
	}

	d.SetId(*read.ID)

	return resourceArmExpressRouteCircuitPeeringRead(d, meta)
}

func resourceArmExpressRouteCircuitPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	peeringType := id.Path["peerings"]

	resp, err := client.Get(ctx, resourceGroup, circuitName, peeringType)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Express Route Circuit Peering %q (Circuit %q / Resource Group %q): %+v", peeringType, circuitName, resourceGroup, err)
	}

	d.Set("peering_type", peeringType)
	d.Set("express_route_circuit_name", circuitName)
	d.Set("resource_group_name", resourceGroup)

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

func resourceArmExpressRouteCircuitPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	peeringType := id.Path["peerings"]

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, resourceGroup, circuitName, peeringType)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Express Route Circuit Peering %q (Circuit %q / Resource Group %q): %+v", peeringType, circuitName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for Express Route Circuit Peering %q (Circuit %q / Resource Group %q) to be deleted: %+v", peeringType, circuitName, resourceGroup, err)
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
		if _, err := ParseRouteFilterID(routeFilterId); err != nil {
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
