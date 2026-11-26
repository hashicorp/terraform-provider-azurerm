// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmLoadBalancerBackendAddressPoolFlatten(d *pluginsdk.ResourceData, id *loadbalancers.LoadBalancerBackendAddressPoolId, model *loadbalancers.BackendAddressPool) error {
	d.Set("name", id.BackendAddressPoolName)
	lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	d.Set("loadbalancer_id", lbId.ID())

	if model != nil {
		if properties := model.Properties; properties != nil {
			if err := d.Set("tunnel_interface", flattenGatewayLoadBalancerTunnelInterfaces(properties.TunnelInterfaces)); err != nil {
				return fmt.Errorf("setting `tunnel_interface`: %v", err)
			}

			var backendIPConfigurations []string
			if configs := properties.BackendIPConfigurations; configs != nil {
				for _, backendConfig := range *configs {
					if backendConfig.Id == nil {
						continue
					}
					backendIPConfigurations = append(backendIPConfigurations, *backendConfig.Id)
				}
			}
			if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
				return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
			}

			d.Set("synchronous_mode", pointer.From(properties.SyncMode))

			network := ""
			if vnet := properties.VirtualNetwork; vnet != nil && vnet.Id != nil {
				network = *vnet.Id
			}
			d.Set("virtual_network_id", network)

			var loadBalancingRules []string
			if rules := properties.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					loadBalancingRules = append(loadBalancingRules, *rule.Id)
				}
			}
			if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
				return fmt.Errorf("setting `load_balancing_rules`: %v", err)
			}

			var outboundRules []string
			if rules := properties.OutboundRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					outboundRules = append(outboundRules, *rule.Id)
				}
			}
			if err := d.Set("outbound_rules", outboundRules); err != nil {
				return fmt.Errorf("setting `outbound_rules`: %v", err)
			}

			var inboundNATRules []string
			if rules := properties.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					inboundNATRules = append(inboundNATRules, *rule.Id)
				}
			}
			if err := d.Set("inbound_nat_rules", inboundNATRules); err != nil {
				return fmt.Errorf("setting `inbound_nat_rules`: %v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenGatewayLoadBalancerTunnelInterfaces(input *[]loadbalancers.GatewayLoadBalancerTunnelInterface) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)

	for _, e := range *input {
		var identifier int
		if e.Identifier != nil {
			identifier = int(*e.Identifier)
		}

		t := pointer.FromEnum(e.Type)

		protocol := pointer.FromEnum(e.Protocol)

		var port int
		if e.Port != nil {
			port = int(*e.Port)
		}

		output = append(output, map[string]any{
			"identifier": identifier,
			"type":       t,
			"protocol":   protocol,
			"port":       port,
		})
	}

	return output
}
