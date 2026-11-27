// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbNatRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NAT Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"loadbalancer_id": {
			Description:  "The ID of the Load Balancer in which to create the NAT Rule.",
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

		"frontend_port": {
			Description:   "The port for the external endpoint. Port numbers for each Rule must be unique within the Load Balancer.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumberOrZero,
			ConflictsWith: []string{"frontend_port_start", "frontend_port_end", "backend_address_pool_id"},
		},

		"backend_port": {
			Description:  "The port used for internal connections on the endpoint.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
		},

		"frontend_ip_configuration_name": {
			Description:  "The name of the frontend IP configuration exposing this rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"floating_ip_enabled": {
			Description: "Whether the Floating IPs are enabled for this Load Balancer Rule. A floating IP is reassigned to a secondary server in case the primary server fails. Required to configure a SQL AlwaysOn Availability Group.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"tcp_reset_enabled": {
			Description: "Whether TCP Reset is enabled for this Load Balancer Rule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"backend_address_pool_id": {
			Description:   "A reference to backendAddressPool resource.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
			RequiredWith:  []string{"frontend_port_start", "frontend_port_end"},
			ConflictsWith: []string{"frontend_port"},
		},

		"frontend_port_start": {
			Description:   "The port range start for the external endpoint. This property is used together with BackendAddressPool and FrontendPortRangeEnd. Individual inbound NAT rule port mappings will be created for each backend address from BackendAddressPool. Acceptable values range from 1 to 65534, inclusive.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumber,
			RequiredWith:  []string{"backend_address_pool_id", "frontend_port_end"},
			ConflictsWith: []string{"frontend_port"},
		},

		"frontend_port_end": {
			Description:   "The port range end for the external endpoint. This property is used together with BackendAddressPool and FrontendPortRangeStart. Individual inbound NAT rule port mappings will be created for each backend address from BackendAddressPool. Acceptable values range from 1 to 65534, inclusive.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumber,
			RequiredWith:  []string{"backend_address_pool_id", "frontend_port_start"},
			ConflictsWith: []string{"frontend_port"},
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

		"backend_ip_configuration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
