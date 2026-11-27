// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	loadbalancerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Load Balancer.",
		Read:        dataSourceArmLoadBalancerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Load Balancer.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku": {
				Description: "The SKU of the Load Balancer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"frontend_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Frontend IP Configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"subnet_id": {
							Description: "The ID of the Subnet which is associated with the IP Configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"private_ip_address": {
							Description: "Private IP Address to assign to the Load Balancer.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"private_ip_address_version": {
							Description: "The Private IP Address Version.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"public_ip_address_id": {
							Description: "The ID of a Public IP Address which is associated with this Load Balancer.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"private_ip_address_allocation": {
							Description: "The allocation method for the Private IP Address used by this Load Balancer.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"zones": commonschema.ZonesMultipleComputed(),

						"id": {
							Description: "The ID of the Frontend IP Configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"private_ip_address": {
				Description: "The first private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"private_ip_addresses": {
				Description: "The list of private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmLoadBalancerRead(d *pluginsdk.ResourceData, meta any) error {
	client := loadbalancerClient.From(meta).LoadBalancersClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
			d.Set("sku", pointer.FromEnum(sku.Name))
		}

		privateIpAddress := ""
		privateIpAddresses := make([]string, 0)
		frontendIpConfigurations := make([]any, 0)

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

func flattenLoadBalancerDataSourceFrontendIpConfiguration(ipConfigs *[]loadbalancers.FrontendIPConfiguration) []any {
	result := make([]any, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		name := pointer.From(config.Name)

		id := pointer.From(config.Id)

		privateIpAddress := ""
		privateIpAddressAllocation := ""
		privateIpAddressVersion := ""
		publicIpAddressId := ""
		subnetId := ""
		if props := config.Properties; props != nil {
			privateIpAddressAllocation = pointer.FromEnum(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				subnetId = pointer.From(subnet.Id)
			}

			privateIpAddress = pointer.From(props.PrivateIPAddress)
			privateIpAddressVersion = pointer.FromEnum(props.PrivateIPAddressVersion)

			if pip := props.PublicIPAddress; pip != nil {
				publicIpAddressId = pointer.From(pip.Id)
			}
		}

		result = append(result, map[string]any{
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
