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
			Description:  "The name of the Outbound Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"loadbalancer_id": {
			Description:  "The ID of the Load Balancer in which to create the Outbound Rule.",
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
						Description:  "The name of the Frontend IP Configuration.",
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
			Description: "The ID of the Backend Address Pool. Outbound traffic is randomly load balanced across IPs in the backend IPs.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"protocol": {
			Description: "The transport protocol for the external endpoint.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.TransportProtocolAll),
				string(loadbalancers.TransportProtocolTcp),
				string(loadbalancers.TransportProtocolUdp),
			}, false),
		},

		"tcp_reset_enabled": {
			Description: "Receive bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination. This element is only used when the protocol is set to TCP.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"allocated_outbound_ports": {
			Description:  "The number of outbound ports to be used for NAT.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1024,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"idle_timeout_in_minutes": {
			Description: "The timeout for the TCP idle connection.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     4,
		},
	}
}
