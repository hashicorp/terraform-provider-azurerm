// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerNatPoolRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatPoolID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
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
		config, _, exists := FindLoadBalancerNatPoolByName(model, id.InboundNatPoolName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Nat Pool %q not found. Removing from state", id.InboundNatPoolName)
			return nil
		}

		d.Set("name", config.Name)
		d.Set("resource_group_name", id.ResourceGroup)

		if props := config.Properties; props != nil {
			d.Set("backend_port", props.BackendPort)
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
			d.Set("frontend_ip_configuration_id", frontendIPConfigID)
			d.Set("frontend_ip_configuration_name", frontendIPConfigName)
			d.Set("frontend_port_end", props.FrontendPortRangeEnd)
			d.Set("frontend_port_start", props.FrontendPortRangeStart)
			d.Set("idle_timeout_in_minutes", int(*props.IdleTimeoutInMinutes))
			d.Set("protocol", string(props.Protocol))
		}
	}

	return nil
}
