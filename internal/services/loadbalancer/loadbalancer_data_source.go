// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

	id := loadbalancers.NewLoadBalancerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: subscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	resp, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if sku := model.Sku; sku != nil {
			d.Set("sku", string(pointer.From(sku.Name)))
		}

		privateIpAddress := ""
		privateIpAddresses := make([]string, 0)
		frontendIpConfigurations := make([]interface{}, 0)

		if props := model.Properties; props != nil {
			if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
				frontendIpConfigurations = flattenLoadBalancerDataSourceFrontendIpConfiguration(feipConfigs)
				for _, config := range *feipConfigs {
					if feipProps := config.Properties; feipProps != nil {
						if ip := feipProps.PrivateIPAddress; ip != nil {
							if privateIpAddress == "" {
								privateIpAddress = pointer.From(ip)
							}

							privateIpAddresses = append(privateIpAddresses, pointer.From(ip))
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

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func flattenLoadBalancerDataSourceFrontendIpConfiguration(ipConfigs *[]loadbalancers.FrontendIPConfiguration) []interface{} {
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
		if config.Id != nil {
			id = *config.Id
		}

		privateIpAddress := ""
		privateIpAddressAllocation := ""
		privateIpAddressVersion := ""
		publicIpAddressId := ""
		subnetId := ""
		if props := config.Properties; props != nil {
			privateIpAddressAllocation = string(pointer.From(props.PrivateIPAllocationMethod))

			if subnet := props.Subnet; subnet != nil {
				subnetId = pointer.From(subnet.Id)
			}

			privateIpAddress = pointer.From(props.PrivateIPAddress)
			privateIpAddressVersion = string(pointer.From(props.PrivateIPAddressVersion))

			if pip := props.PublicIPAddress; pip != nil {
				publicIpAddressId = pointer.From(pip.Id)
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
