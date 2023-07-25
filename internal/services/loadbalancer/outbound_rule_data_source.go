// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceArmLoadBalancerOutboundRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmLoadBalancerOutboundRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.LoadBalancerID,
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

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backend_address_pool_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tcp_reset_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"allocated_outbound_ports": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceArmLoadBalancerOutboundRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	loadBalancer, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			return fmt.Errorf("parent %s was not found", *loadBalancerId)
		}
		return fmt.Errorf("retrieving parent %s: %+v", *loadBalancerId, err)
	}

	id := parse.NewLoadBalancerOutboundRuleID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, name)
	config, _, exists := FindLoadBalancerOutboundRuleByName(&loadBalancer, id.OutboundRuleName)
	if !exists {
		return fmt.Errorf("%s was not found", id)
	}

	d.SetId(id.ID())
	if props := config.OutboundRulePropertiesFormat; props != nil {
		allocatedOutboundPorts := 0
		if props.AllocatedOutboundPorts != nil {
			allocatedOutboundPorts = int(*props.AllocatedOutboundPorts)
		}
		d.Set("allocated_outbound_ports", allocatedOutboundPorts)

		backendAddressPoolId := ""
		if props.BackendAddressPool != nil && props.BackendAddressPool.ID != nil {
			bapid, err := parse.LoadBalancerBackendAddressPoolID(*props.BackendAddressPool.ID)
			if err != nil {
				return err
			}

			backendAddressPoolId = bapid.ID()
		}
		d.Set("backend_address_pool_id", backendAddressPoolId)
		d.Set("tcp_reset_enabled", props.EnableTCPReset)

		frontendIpConfigurations := make([]interface{}, 0)
		if configs := props.FrontendIPConfigurations; configs != nil {
			for _, feConfig := range *configs {
				if feConfig.ID == nil {
					continue
				}
				feid, err := parse.LoadBalancerFrontendIpConfigurationID(*feConfig.ID)
				if err != nil {
					return err
				}

				frontendIpConfigurations = append(frontendIpConfigurations, map[string]interface{}{
					"id":   feid.ID(),
					"name": feid.FrontendIPConfigurationName,
				})
			}
		}
		d.Set("frontend_ip_configuration", frontendIpConfigurations)

		idleTimeoutInMinutes := 0
		if props.IdleTimeoutInMinutes != nil {
			idleTimeoutInMinutes = int(*props.IdleTimeoutInMinutes)
		}
		d.Set("idle_timeout_in_minutes", idleTimeoutInMinutes)
		d.Set("protocol", string(props.Protocol))
	}

	return nil
}
