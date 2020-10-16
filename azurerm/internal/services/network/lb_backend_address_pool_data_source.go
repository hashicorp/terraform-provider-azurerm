package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmLoadBalancerBackendAddressPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLoadBalancerBackendAddressPoolRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.LoadBalancerID,
			},

			"backend_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"virtual_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"network_interface_ip_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backend_ip_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancing_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func dataSourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LoadBalancerBackendAddressPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Load Balancer Backend Address Pool %q (Resource Group %q / Load Balancer %q) was not found", name, loadBalancerId.ResourceGroup, loadBalancerId.Name)
		}

		return fmt.Errorf("Error making Read request on Load Balancer Backend Address Pool %q (Resource Group %q / Load Balancer %q)", name, loadBalancerId.ResourceGroup, loadBalancerId.Name)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("nil or empty ID of Load Balancer Backend Address Pool %q (Resource Group %q / Load Balancer %q)", name, loadBalancerId.ResourceGroup, loadBalancerId.Name)
	}

	id, err := parse.LoadBalancerBackendAddressPoolID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	var backendIpConfigurations []interface{}
	var loadBalancingRules []string
	var outboundRules []string

	if props := resp.BackendAddressPoolPropertiesFormat; props != nil {
		if err := d.Set("backend_address", flattenArmLoadBalancerBackendAddressesForDataSource(props.LoadBalancerBackendAddresses)); err != nil {
			return fmt.Errorf("setting `backend_address`: %v", err)
		}

		if configs := props.BackendIPConfigurations; configs != nil {
			for _, backendConfig := range *configs {
				if backendConfig.ID == nil {
					continue
				}
				v := map[string]interface{}{
					"id": *backendConfig.ID,
				}
				backendIpConfigurations = append(backendIpConfigurations, v)
			}
		}

		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				loadBalancingRules = append(loadBalancingRules, *rule.ID)
			}
		}

		if rules := props.OutboundRules; rules != nil {
			for _, rule := range *rules {
				outboundRules = append(outboundRules, *rule.ID)
			}
		}
	}

	d.Set("backend_ip_configurations", backendIpConfigurations)
	d.Set("load_balancing_rules", loadBalancingRules)
	d.Set("outbound_rules", outboundRules)

	return nil
}

func flattenArmLoadBalancerBackendAddressesForDataSource(input *[]network.LoadBalancerBackendAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var name string
		if e.Name != nil {
			name = *e.Name
		}

		var (
			ipAddress string
			vnetId    string
			ipConfig  string
		)
		if prop := e.LoadBalancerBackendAddressPropertiesFormat; prop != nil {
			if prop.IPAddress != nil {
				ipAddress = *prop.IPAddress
			}
			if prop.VirtualNetwork != nil && prop.VirtualNetwork.ID != nil {
				vnetId = *prop.VirtualNetwork.ID
			}
			if prop.NetworkInterfaceIPConfiguration != nil && prop.NetworkInterfaceIPConfiguration.ID != nil {
				ipConfig = *prop.NetworkInterfaceIPConfiguration.ID
			}
		}

		v := map[string]interface{}{
			"name":                               name,
			"virtual_network_id":                 vnetId,
			"ip_address":                         ipAddress,
			"network_interface_ip_configuration": ipConfig,
		}
		output = append(output, v)
	}

	return output
}
