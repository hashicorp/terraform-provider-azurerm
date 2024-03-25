// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
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
	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: loadBalancerId.SubscriptionId, ResourceGroupName: loadBalancerId.ResourceGroupName, LoadBalancerName: loadBalancerId.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			return fmt.Errorf("%s was not found", *loadBalancerId)
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	id := loadbalancers.NewOutboundRuleID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)

	if model := loadBalancer.Model; model != nil {
		config, _, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
		if !exists {
			return fmt.Errorf("%s was not found", id)
		}

		d.SetId(id.ID())

		if props := config.Properties; props != nil {
			d.Set("allocated_outbound_ports", int(pointer.From(props.AllocatedOutboundPorts)))

			backendAddressPoolId := ""
			if props.BackendAddressPool.Id != nil {
				bapid, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(*props.BackendAddressPool.Id)
				if err != nil {
					return err
				}

				backendAddressPoolId = bapid.ID()
			}
			d.Set("backend_address_pool_id", backendAddressPoolId)
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))

			frontendIpConfigurations := make([]interface{}, 0)
			if configs := props.FrontendIPConfigurations; configs != nil {
				for _, feConfig := range configs {
					if feConfig.Id == nil {
						continue
					}
					feid, err := loadbalancers.ParseFrontendIPConfigurationID(*feConfig.Id)
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
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("protocol", string(props.Protocol))
		}
	}
	return nil
}
