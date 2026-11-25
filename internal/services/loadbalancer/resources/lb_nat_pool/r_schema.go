// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbNatPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NAT pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"loadbalancer_id": {
			Description:  "The ID of the Load Balancer in which to create the NAT pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
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

		"frontend_port_start": {
			Description:  "The first port number in the range of external ports that will be used to provide Inbound NAT to NICs associated with this Load Balancer.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"frontend_port_end": {
			Description:  "The last port number in the range of external ports that will be used to provide Inbound NAT to NICs associated with this Load Balancer.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"backend_port": {
			Description:  "The port used for the internal endpoint.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"frontend_ip_configuration_name": {
			Description:  "The name of the frontend IP configuration exposing this rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"floating_ip_enabled": {
			Description: "Whether the floating IPs are enabled for this Load Balancer Rule. A floating IP is reassigned to a secondary server in case the primary server fails. Required to configure a SQL AlwaysOn Availability Group.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tcp_reset_enabled": {
			Description: "Whether TCP Reset is enabled for this Load Balancer Rule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"idle_timeout_in_minutes": {
			Description:  "The idle timeout in minutes for TCP connections.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"frontend_ip_configuration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
