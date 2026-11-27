// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_rule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmLoadBalancerRuleFlatten(d *pluginsdk.ResourceData, id *loadbalancers.LoadBalancingRuleId, model *loadbalancers.LoadBalancingRule) error {
	d.Set("name", model.Name)
	lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	d.Set("loadbalancer_id", lbId.ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("disable_outbound_snat", pointer.From(props.DisableOutboundSnat))
			d.Set("floating_ip_enabled", pointer.From(props.EnableFloatingIP))
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))
			d.Set("protocol", string(props.Protocol))
			d.Set("backend_port", int(pointer.From(props.BackendPort)))

			var backendAddressPoolIds []string
			if props.BackendAddressPools != nil {
				for _, p := range *props.BackendAddressPools {
					backendAddressPoolIds = append(backendAddressPoolIds, pointer.From(p.Id))
				}
			} else if props.BackendAddressPool != nil {
				backendAddressPoolIds = []string{pointer.From(props.BackendAddressPool.Id)}
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
	return pluginsdk.SetResourceIdentityData(d, id)
}
