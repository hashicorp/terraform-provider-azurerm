package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

						"zones": azure.SchemaZonesComputed(),

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

			"tags": tags.SchemaDataSource(),
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
	if props := resp.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			if err := d.Set("frontend_ip_configuration", flattenLoadBalancerDataSourceFrontendIpConfiguration(feipConfigs)); err != nil {
				return fmt.Errorf("flattening `frontend_ip_configuration`: %+v", err)
			}

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
		ipConfig := make(map[string]interface{})
		if config.Name != nil {
			ipConfig["name"] = *config.Name
		}

		if config.ID != nil {
			ipConfig["id"] = *config.ID
		}

		zones := make([]string, 0)
		if zs := config.Zones; zs != nil {
			zones = *zs
		}
		ipConfig["zones"] = zones

		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			ipConfig["private_ip_address_allocation"] = props.PrivateIPAllocationMethod

			if subnet := props.Subnet; subnet != nil && subnet.ID != nil {
				ipConfig["subnet_id"] = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				ipConfig["private_ip_address"] = *pip
			}

			if props.PrivateIPAddressVersion != "" {
				ipConfig["private_ip_address_version"] = string(props.PrivateIPAddressVersion)
			}

			if pip := props.PublicIPAddress; pip != nil && pip.ID != nil {
				ipConfig["public_ip_address_id"] = *pip.ID
			}
		}

		result = append(result, ipConfig)
	}
	return result
}
