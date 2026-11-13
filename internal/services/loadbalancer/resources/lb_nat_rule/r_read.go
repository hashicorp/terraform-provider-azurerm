// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_rule

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

func resourceArmLoadBalancerNatRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseInboundNatRuleID(d.Id())
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
		config, _, exists := FindLoadBalancerNatRuleByName(model, id.InboundNatRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Nat Rule %q not found. Removing from state", id.InboundNatRuleName)
			return nil
		}

		d.Set("name", config.Name)
		d.Set("resource_group_name", id.ResourceGroupName)

		if props := config.Properties; props != nil {
			backendIPConfigId := ""
			if props.BackendIPConfiguration != nil && props.BackendIPConfiguration.Id != nil {
				backendIPConfigId = *props.BackendIPConfiguration.Id
			}
			d.Set("backend_ip_configuration_id", backendIPConfigId)
			d.Set("backend_port", pointer.From(props.BackendPort))
			d.Set("floating_ip_enabled", pointer.From(props.EnableFloatingIP))
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))

			frontendIPConfigName := ""
			frontendIPConfigID := ""
			if props.FrontendIPConfiguration != nil && props.FrontendIPConfiguration.Id != nil {
				feid, err := loadbalancers.ParseFrontendIPConfigurationIDInsensitively(*props.FrontendIPConfiguration.Id)
				if err != nil {
					return err
				}

				frontendIPConfigName = feid.FrontendIPConfigurationName
				frontendIPConfigID = feid.ID()
			}
			d.Set("frontend_ip_configuration_name", frontendIPConfigName)
			d.Set("frontend_ip_configuration_id", frontendIPConfigID)

			if props.BackendAddressPool != nil {
				d.Set("backend_address_pool_id", pointer.From(props.BackendAddressPool.Id))
			}
			d.Set("frontend_port", pointer.From(props.FrontendPort))
			d.Set("frontend_port_start", int(pointer.From(props.FrontendPortRangeStart)))
			d.Set("frontend_port_end", int(pointer.From(props.FrontendPortRangeEnd)))
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("protocol", pointer.FromEnum(props.Protocol))
		}
	}
	return nil
}
