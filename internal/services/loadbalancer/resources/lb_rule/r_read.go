// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_rule

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

func resourceArmLoadBalancerRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancingRuleID(d.Id())
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
		config, _, exists := FindLoadBalancerRuleByName(model, id.LoadBalancingRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Rule %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}

		d.Set("name", config.Name)

		if props := config.Properties; props != nil {
			d.Set("disable_outbound_snat", pointer.From(props.DisableOutboundSnat))
			d.Set("floating_ip_enabled", pointer.From(props.EnableFloatingIP))
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))
			d.Set("protocol", string(props.Protocol))
			d.Set("backend_port", int(pointer.From(props.BackendPort)))

			// The backendAddressPools is designed for Gateway LB, while the backendAddressPool is designed for other skus.
			// Thought currently the API returns both, but for the sake of stability, we do use different fields here depending on the LB sku.
			var isGateway bool
			if model.Sku != nil && pointer.From(model.Sku.Name) == loadbalancers.LoadBalancerSkuNameGateway {
				isGateway = true
			}
			var (
				backendAddressPoolId  string
				backendAddressPoolIds []any
			)
			if isGateway {
				// The gateway LB rule can have up to 2 backend address pools.
				// In case there is only one BAP, we set it to both "backendAddressPoolId" and "backendAddressPoolIds".
				// Otherwise, we leave the "backendAddressPoolId" as empty.
				if props.BackendAddressPools != nil {
					for _, p := range *props.BackendAddressPools {
						if p.Id != nil {
							backendAddressPoolIds = append(backendAddressPoolIds, *p.Id)
						}
					}
				}
			} else {
				if props.BackendAddressPool != nil && props.BackendAddressPool.Id != nil {
					backendAddressPoolId = *props.BackendAddressPool.Id
					backendAddressPoolIds = []any{backendAddressPoolId}
				}
			}
			d.Set("backend_address_pool_ids", backendAddressPoolIds)

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
			d.Set("frontend_port", int(props.FrontendPort))
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("load_distribution", pointer.FromEnum(props.LoadDistribution))

			probeId := ""
			if props.Probe != nil {
				probeId = pointer.From(props.Probe.Id)
			}
			d.Set("probe_id", probeId)
		}
	}
	return nil
}
