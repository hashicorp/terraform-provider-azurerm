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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"backend_address_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
		},

		"virtual_network_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"backend_address_ip_configuration_id"},
			ValidateFunc:  commonids.ValidateVirtualNetworkID,
		},

		"ip_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPAddress,
		},

		"backend_address_ip_configuration_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"virtual_network_id"},
			ValidateFunc:  loadbalancers.ValidateFrontendIPConfigurationID,
			Description:   "For global load balancer, user needs to specify the `backend_address_ip_configuration_id` of the added regional load balancers",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"inbound_nat_rule_port_mapping": portMapping(),
	}
}
