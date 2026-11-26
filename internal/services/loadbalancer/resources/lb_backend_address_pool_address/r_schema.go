// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Backend Address Pool Address.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"backend_address_pool_id": {
			Description:  "The ID of the Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
		},

		"virtual_network_id": {
			Description:   "The ID of the Virtual Network within which the Backend Address Pool should exist.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateVirtualNetworkID,
			ConflictsWith: []string{"backend_address_ip_configuration_id"},
		},

		"ip_address": {
			Description:  "The Static IP Address which should be allocated to this Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPAddress,
		},

		"backend_address_ip_configuration_id": {
			Description:   "The IP config ID of the regional load balancer that's added to the global load balancer's backend address pool.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  loadbalancers.ValidateFrontendIPConfigurationID,
			ConflictsWith: []string{"virtual_network_id"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"inbound_nat_rule_port_mapping": portMapping(),
	}
}
