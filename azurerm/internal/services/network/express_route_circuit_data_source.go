package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceExpressRouteCircuit() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceExpressRouteCircuitRead,

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

			"peerings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"peering_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"primary_peer_address_prefix": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"secondary_peer_address_prefix": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"azure_asn": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"peer_asn": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"vlan_id": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"shared_key": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"service_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_provider_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"service_provider_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"peering_location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"bandwidth_in_mbps": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"service_provider_provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceExpressRouteCircuitRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Express Route Circuit %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on the Express Route Circuit %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.ExpressRouteCircuitPropertiesFormat; properties != nil {
		peerings := flattenExpressRouteCircuitPeerings(properties.Peerings)
		if err := d.Set("peerings", peerings); err != nil {
			return err
		}

		d.Set("service_key", properties.ServiceKey)
		d.Set("service_provider_provisioning_state", properties.ServiceProviderProvisioningState)

		if serviceProviderProperties := flattenExpressRouteCircuitServiceProviderProperties(properties.ServiceProviderProperties); serviceProviderProperties != nil {
			if err := d.Set("service_provider_properties", serviceProviderProperties); err != nil {
				return fmt.Errorf("Error setting `service_provider_properties`: %+v", err)
			}
		}
	}

	if resp.Sku != nil {
		sku := flattenExpressRouteCircuitSku(resp.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("Error setting `sku`: %+v", err)
		}
	}

	return nil
}

func flattenExpressRouteCircuitPeerings(input *[]network.ExpressRouteCircuitPeering) []interface{} {
	peerings := make([]interface{}, 0)

	if input != nil {
		for _, peering := range *input {
			props := peering.ExpressRouteCircuitPeeringPropertiesFormat
			p := make(map[string]interface{})

			p["peering_type"] = string(props.PeeringType)

			if primaryPeerAddressPrefix := props.PrimaryPeerAddressPrefix; primaryPeerAddressPrefix != nil {
				p["primary_peer_address_prefix"] = *primaryPeerAddressPrefix
			}

			if secondaryPeerAddressPrefix := props.SecondaryPeerAddressPrefix; secondaryPeerAddressPrefix != nil {
				p["secondary_peer_address_prefix"] = *secondaryPeerAddressPrefix
			}

			if azureAsn := props.AzureASN; azureAsn != nil {
				p["azure_asn"] = *azureAsn
			}

			if peerAsn := props.PeerASN; peerAsn != nil {
				p["peer_asn"] = *peerAsn
			}

			if vlanID := props.VlanID; vlanID != nil {
				p["vlan_id"] = *vlanID
			}

			if sharedKey := props.SharedKey; sharedKey != nil {
				p["shared_key"] = *sharedKey
			}

			peerings = append(peerings, p)
		}
	}

	return peerings
}

func flattenExpressRouteCircuitServiceProviderProperties(input *network.ExpressRouteCircuitServiceProviderProperties) []interface{} {
	serviceProviderProperties := make([]interface{}, 0)

	if input != nil {
		p := make(map[string]interface{})

		if serviceProviderName := input.ServiceProviderName; serviceProviderName != nil {
			p["service_provider_name"] = *serviceProviderName
		}

		if peeringLocation := input.PeeringLocation; peeringLocation != nil {
			p["peering_location"] = *peeringLocation
		}

		if bandwidthInMbps := input.BandwidthInMbps; bandwidthInMbps != nil {
			p["bandwidth_in_mbps"] = *bandwidthInMbps
		}

		serviceProviderProperties = append(serviceProviderProperties, p)
	}

	return serviceProviderProperties
}
