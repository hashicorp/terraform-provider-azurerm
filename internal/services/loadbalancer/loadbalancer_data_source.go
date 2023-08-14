// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func dataSourceArmLoadBalancer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmLoadBalancerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"frontend_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"zones": commonschema.ZonesMultipleComputed(),

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmLoadBalancerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLoadBalancerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("location", location.NormalizeNilable(resp.Location))
	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	privateIpAddress := ""
	privateIpAddresses := make([]string, 0)
	frontendIpConfigurations := make([]interface{}, 0)

	if props := resp.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			frontendIpConfigurations = flattenLoadBalancerDataSourceFrontendIpConfiguration(feipConfigs)

			for _, config := range *feipConfigs {
				if feipProps := config.FrontendIPConfigurationPropertiesFormat; feipProps != nil {
					if ip := feipProps.PrivateIPAddress; ip != nil {
						if privateIpAddress == "" {
							privateIpAddress = *feipProps.PrivateIPAddress
						}

						privateIpAddresses = append(privateIpAddresses, *feipProps.PrivateIPAddress)
					}
				}
			}
		}
	}

	if err := d.Set("frontend_ip_configuration", frontendIpConfigurations); err != nil {
		return fmt.Errorf("flattening `frontend_ip_configuration`: %+v", err)
	}
	d.Set("private_ip_address", privateIpAddress)
	d.Set("private_ip_addresses", privateIpAddresses)

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenLoadBalancerDataSourceFrontendIpConfiguration(ipConfigs *[]network.FrontendIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		name := ""
		if config.Name != nil {
			name = *config.Name
		}

		id := ""
		if config.ID != nil {
			id = *config.ID
		}

		privateIpAddress := ""
		privateIpAddressAllocation := ""
		privateIpAddressVersion := ""
		publicIpAddressId := ""
		subnetId := ""
		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			privateIpAddressAllocation = string(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil && subnet.ID != nil {
				subnetId = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				privateIpAddress = *pip
			}

			if props.PrivateIPAddressVersion != "" {
				privateIpAddressVersion = string(props.PrivateIPAddressVersion)
			}

			if pip := props.PublicIPAddress; pip != nil && pip.ID != nil {
				publicIpAddressId = *pip.ID
			}
		}

		result = append(result, map[string]interface{}{
			"id":                            id,
			"name":                          name,
			"private_ip_address":            privateIpAddress,
			"private_ip_address_allocation": privateIpAddressAllocation,
			"private_ip_address_version":    privateIpAddressVersion,
			"public_ip_address_id":          publicIpAddressId,
			"subnet_id":                     subnetId,
			"zones":                         zones.FlattenUntyped(config.Zones),
		})
	}
	return result
}
