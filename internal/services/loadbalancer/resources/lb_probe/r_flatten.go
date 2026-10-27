// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_probe

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceArmLoadBalancerProbeFlatten(d *pluginsdk.ResourceData, id *loadbalancers.ProbeId, model *loadbalancers.Probe) error {
	d.Set("name", model.Name)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("interval_in_seconds", int(pointer.From(props.IntervalInSeconds)))
			d.Set("number_of_probes", int(pointer.From(props.NumberOfProbes)))
			d.Set("port", int(props.Port))
			d.Set("protocol", string(props.Protocol))
			d.Set("request_path", pointer.From(props.RequestPath))
			d.Set("probe_threshold", int(pointer.From(props.ProbeThreshold)))
			d.Set("no_healthy_backends_behavior", pointer.FromEnum(props.NoHealthyBackendsBehavior))

			// TODO: parse/make these consistent
			var loadBalancerRules []string
			if rules := props.LoadBalancingRules; rules != nil {
				for _, ruleConfig := range *rules {
					if id := ruleConfig.Id; id != nil {
						loadBalancerRules = append(loadBalancerRules, *id)
					}
				}
			}
			if err := d.Set("load_balancer_rules", loadBalancerRules); err != nil {
				return fmt.Errorf("setting `load_balancer_rules` (Load Balancer %q Probe %q): %+v", id.LoadBalancerName, id.ProbeName, err)
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
