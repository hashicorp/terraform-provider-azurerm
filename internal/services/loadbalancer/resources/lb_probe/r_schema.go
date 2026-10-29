// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_probe

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbProbeSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"loadbalancer_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
		},

		"protocol": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.ProbeProtocolTcp),
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForProbeProtocol(), false),
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"probe_threshold": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"request_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"interval_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      15,
			ValidateFunc: validation.IntAtLeast(5),
		},

		"number_of_probes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  2,
		},

		"no_healthy_backends_behavior": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.ProbeNoHealthyBackendsBehaviorAllProbedDown),
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForProbeNoHealthyBackendsBehavior(), false),
		},

		"load_balancer_rules": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},
	}
}
