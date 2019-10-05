package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
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
		},
	}
}

func resourceArmExpressRouteCircuitPeeringCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ExpressRoutePeeringsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Express Route Peering create/update.")

	peeringType := d.Get("peering_type").(string)
	circuitName := d.Get("express_route_circuit_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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
	client := meta.(*ArmClient).network.ExpressRoutePeeringsClient
	ctx := meta.(*ArmClient).StopContext

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

		config := flattenExpressRouteCircuitPeeringMicrosoftConfig(props.MicrosoftPeeringConfig)
		if err := d.Set("microsoft_peering_config", config); err != nil {
			return fmt.Errorf("Error setting `microsoft_peering_config`: %+v", err)
		}
	}

	return nil
}

func resourceArmExpressRouteCircuitPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ExpressRoutePeeringsClient
	ctx := meta.(*ArmClient).StopContext

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
	peering := input[0].(map[string]interface{})

	prefixes := make([]string, 0)
	inputPrefixes := peering["advertised_public_prefixes"].([]interface{})
	for _, v := range inputPrefixes {
		prefixes = append(prefixes, v.(string))
	}

	return &network.ExpressRouteCircuitPeeringConfig{
		AdvertisedPublicPrefixes: &prefixes,
	}
}

func flattenExpressRouteCircuitPeeringMicrosoftConfig(input *network.ExpressRouteCircuitPeeringConfig) interface{} {
	if input == nil {
		return []interface{}{}
	}

	config := make(map[string]interface{})
	prefixes := make([]string, 0)
	if ps := input.AdvertisedPublicPrefixes; ps != nil {
		prefixes = *ps
	}

	config["advertised_public_prefixes"] = prefixes

	return []interface{}{config}
}
