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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	client := meta.(*clients.Client).Network.ExpressRouteCircuits
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := expressroutecircuits.NewExpressRouteCircuitID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		sku := flattenExpressRouteCircuitSku(model.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}
		if props := model.Properties; props != nil {
			peerings := flattenExpressRouteCircuitPeerings(props.Peerings)
			if err := d.Set("peerings", peerings); err != nil {
				return err
			}

			d.Set("service_key", props.ServiceKey)
			d.Set("service_provider_provisioning_state", string(pointer.From(props.ServiceProviderProvisioningState)))

			if serviceProviderProperties := flattenExpressRouteCircuitServiceProviderProperties(props.ServiceProviderProperties); serviceProviderProperties != nil {
				if err := d.Set("service_provider_properties", serviceProviderProperties); err != nil {
					return fmt.Errorf("setting `service_provider_properties`: %+v", err)
				}
			}
		}
	}

	return nil
}

func flattenExpressRouteCircuitPeerings(input *[]expressroutecircuits.ExpressRouteCircuitPeering) []interface{} {
	peerings := make([]interface{}, 0)

	if input != nil {
		for _, peering := range *input {
			props := peering.Properties
			p := make(map[string]interface{})

			p["peering_type"] = string(pointer.From(props.PeeringType))

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

			if vlanID := props.VlanId; vlanID != nil {
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

func flattenExpressRouteCircuitServiceProviderProperties(input *expressroutecircuits.ExpressRouteCircuitServiceProviderProperties) []interface{} {
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
