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
			Description:  "The name of the Probe.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"loadbalancer_id": {
			Description:  "The ID of the LoadBalancer in which to create the Probe.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
		},

		"protocol": {
			Description:  "The protocol of the end point. If TCP is specified, a received ACK is required for the probe to be successful. If HTTP is specified, a 200 OK response from the specified URI is required for the probe to be successful.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.ProbeProtocolTcp),
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForProbeProtocol(), false),
		},

		"port": {
			Description:  "Port on which the Probe queries the backend endpoint.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"probe_threshold": {
			Description:  "The number of consecutive successful or failed probes that allow or deny traffic to this endpoint.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"request_path": {
			Description: "The URI used for requesting health status from the backend endpoint. Otherwise, it is not allowed.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"interval_in_seconds": {
			Description:  "The interval, in seconds between probes to the backend endpoint for health status.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      15,
			ValidateFunc: validation.IntAtLeast(5),
		},

		"number_of_probes": {
			Description: "The number of failed probe attempts after which the backend endpoint is removed from rotation. NumberOfProbes multiplied by intervalInSeconds value must be greater or equal to 10. Endpoints are returned to rotation when at least one probe is successful.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     2,
		},

		"no_healthy_backends_behavior": {
			Description:  "How new connections are handled by the load balancer when all backend instances are probed down.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.ProbeNoHealthyBackendsBehaviorAllProbedDown),
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForProbeNoHealthyBackendsBehavior(), false),
		},

		"load_balancer_rules": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
