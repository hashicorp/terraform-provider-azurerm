// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbBackendAddressPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"loadbalancer_id": {
			Description:  "The ID of the Load Balancer in which to create the Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
		},

		"synchronous_mode": {
			Description:  "The backend address synchronous mode for the Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForSyncMode(), false),
			RequiredWith: []string{"virtual_network_id"},
		},

		"tunnel_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identifier": {
						Description: "The unique identifier of this Gateway Load Balancer Tunnel Interface.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"type": {
						Description: "The traffic type of this Gateway Load Balancer Tunnel Interface.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice(
							loadbalancers.PossibleValuesForGatewayLoadBalancerTunnelInterfaceType(),
							false,
						),
					},

					"protocol": {
						Description: "The protocol used for this Gateway Load Balancer Tunnel Interface.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice(
							loadbalancers.PossibleValuesForGatewayLoadBalancerTunnelProtocol(),
							false,
						),
					},

					"port": {
						Description: "The port number that this Gateway Load Balancer Tunnel Interface listens to.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},
				},
			},
		},

		"virtual_network_id": {
			Description:  "The ID of the Virtual Network within which the Backend Address Pool should exist.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"backend_ip_configurations": {
			Description: "The Backend IP Configurations associated with this Backend Address Pool.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"inbound_nat_rules": {
			Description: "An array of the Load Balancing Inbound NAT Rules associated with this Backend Address Pool.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"load_balancing_rules": {
			Description: "The Load Balancing Rules associated with this Backend Address Pool.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"outbound_rules": {
			Description: "An array of the Load Balancing Outbound Rules associated with this Backend Address Pool.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
