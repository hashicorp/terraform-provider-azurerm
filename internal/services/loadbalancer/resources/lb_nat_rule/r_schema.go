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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"loadbalancer_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
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

		"frontend_port": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumberOrZero,
			ConflictsWith: []string{"frontend_port_start", "frontend_port_end", "backend_address_pool_id"},
		},

		"backend_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
		},

		"frontend_ip_configuration_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"floating_ip_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"tcp_reset_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"backend_address_pool_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
			ConflictsWith: []string{"frontend_port"},
			RequiredWith:  []string{"frontend_port_start", "frontend_port_end"},
		},

		"frontend_port_start": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumber,
			RequiredWith:  []string{"backend_address_pool_id", "frontend_port_end"},
			ConflictsWith: []string{"frontend_port"},
		},

		"frontend_port_end": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IsPortNumber,
			RequiredWith:  []string{"backend_address_pool_id", "frontend_port_start"},
			ConflictsWith: []string{"frontend_port"},
		},

		"idle_timeout_in_minutes": {
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
