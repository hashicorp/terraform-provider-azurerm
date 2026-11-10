// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerOutboundRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseOutboundRuleID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", plbId, err)
	}

	if model := loadBalancer.Model; model != nil {
		config, _, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Outbound Rule %q not found. Removing from state", id.OutboundRuleName)
			return nil
		}

		d.Set("name", config.Name)

		if props := config.Properties; props != nil {
			allocatedOutboundPorts := 0
			if props.AllocatedOutboundPorts != nil {
				allocatedOutboundPorts = int(*props.AllocatedOutboundPorts)
			}
			d.Set("allocated_outbound_ports", allocatedOutboundPorts)

			backendAddressPoolId := ""
			if props.BackendAddressPool.Id != nil {
				bapid, err := loadbalancers.ParseLoadBalancerBackendAddressPoolIDInsensitively(*props.BackendAddressPool.Id)
				if err != nil {
					return err
				}

				backendAddressPoolId = bapid.ID()
			}
			d.Set("backend_address_pool_id", backendAddressPoolId)
			d.Set("tcp_reset_enabled", props.EnableTcpReset)

			frontendIpConfigurations := make([]any, 0)
			if configs := props.FrontendIPConfigurations; configs != nil {
				for _, feConfig := range configs {
					if feConfig.Id == nil {
						continue
					}
					feid, err := loadbalancers.ParseFrontendIPConfigurationIDInsensitively(*feConfig.Id)
					if err != nil {
						return err
					}

					frontendIpConfigurations = append(frontendIpConfigurations, map[string]any{
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
