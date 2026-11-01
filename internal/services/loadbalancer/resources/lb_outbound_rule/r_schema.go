// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbOutboundRuleSchema() map[string]*pluginsdk.Schema {
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

		"frontend_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
			Required: true,
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.TransportProtocolAll),
				string(loadbalancers.TransportProtocolTcp),
				string(loadbalancers.TransportProtocolUdp),
			}, false),
		},

		"tcp_reset_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"allocated_outbound_ports": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1024,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"idle_timeout_in_minutes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  4,
		},
	}
}
