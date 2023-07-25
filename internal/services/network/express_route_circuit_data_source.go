// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewExpressRouteCircuitID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("location", location.NormalizeNilable(resp.Location))

	if properties := resp.ExpressRouteCircuitPropertiesFormat; properties != nil {
		peerings := flattenExpressRouteCircuitPeerings(properties.Peerings)
		if err := d.Set("peerings", peerings); err != nil {
			return err
		}

		d.Set("service_key", properties.ServiceKey)
		d.Set("service_provider_provisioning_state", properties.ServiceProviderProvisioningState)

		if serviceProviderProperties := flattenExpressRouteCircuitServiceProviderProperties(properties.ServiceProviderProperties); serviceProviderProperties != nil {
			if err := d.Set("service_provider_properties", serviceProviderProperties); err != nil {
				return fmt.Errorf("setting `service_provider_properties`: %+v", err)
			}
		}
	}

	sku := flattenExpressRouteCircuitSku(resp.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
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
