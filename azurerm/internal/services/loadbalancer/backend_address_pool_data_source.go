package loadbalancer

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
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
					Type: schema.TypeString,
				},
			},

			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmLoadBalancerBackendAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, name)

	resp, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Load Balancer Backend Address Pool %q was not found", id)
		}
		return fmt.Errorf("failed to retrieve Load Balancer Backend Address Pool %q: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := resp.BackendAddressPoolPropertiesFormat; props != nil {
		if err := d.Set("backend_address", flattenArmLoadBalancerBackendAddresses(props.LoadBalancerBackendAddresses)); err != nil {
			return fmt.Errorf("setting `backend_address`: %v", err)
		}

		var backendIPConfigurations []interface{}
		if beipConfigs := props.BackendIPConfigurations; beipConfigs != nil {
			for _, config := range *beipConfigs {
				ipConfig := make(map[string]interface{})
				if id := config.ID; id != nil {
					ipConfig["id"] = *id
					backendIPConfigurations = append(backendIPConfigurations, ipConfig)
				}
			}
		}
		if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
			return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
		}

		var loadBalancingRules []string
		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				loadBalancingRules = append(loadBalancingRules, *rule.ID)
			}
		}
		if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
			return fmt.Errorf("setting `load_balancing_rules`: %v", err)
		}

		var outboundRules []string
		if rules := props.OutboundRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				outboundRules = append(outboundRules, *rule.ID)
			}
		}
		if err := d.Set("outbound_rules", outboundRules); err != nil {
			return fmt.Errorf("setting `outbound_rules`: %v", err)
		}
	}

	return nil
}
