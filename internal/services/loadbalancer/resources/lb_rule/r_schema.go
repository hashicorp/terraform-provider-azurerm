// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceArmLoadBalancerRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the LB Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RuleName,
		},

		"loadbalancer_id": {
			Description:  "The ID of the Load Balancer in which to create the Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
		},

		"frontend_ip_configuration_name": {
			Description:  "The name of the frontend IP configuration to which the rule is associated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frontend_ip_configuration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"backend_address_pool_ids": {
			Description: "A list of reference to a Backend Address Pool over which this Load Balancing Rule operates.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    2, // Only Gateway SKU LB can have 2 backend address pools
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
			},
		},

		"protocol": {
			Description:      "The transport protocol for the external endpoint.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.TransportProtocolAll),
				string(loadbalancers.TransportProtocolTcp),
				string(loadbalancers.TransportProtocolUdp),
			}, false),
		},

		"frontend_port": {
			Description:  "The port for the external endpoint. Port numbers for each Rule must be unique within the Load Balancer. A port of `0` means \"Any Port\".",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
		},

		"backend_port": {
			Description:  "The port used for internal connections on the endpoint. A port of `0` means \"Any Port\".",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
		},

		"probe_id": {
			Description: "A reference to a Probe used by this Load Balancing Rule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"floating_ip_enabled": {
			Description: "Are the Floating IPs enabled for this Load Balancer Rule? A \"floating\" IP is reassigned to a secondary server in case the primary server fails. Required to configure a SQL AlwaysOn Availability Group.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tcp_reset_enabled": {
			Description: "Is TCP Reset enabled for this Load Balancer Rule?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"disable_outbound_snat": {
			Description: "Is snat enabled for this Load Balancer Rule?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"idle_timeout_in_minutes": {
			Description:  "The idle timeout in minutes for TCP connections.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 100),
		},

		"load_distribution": {
			Description: "The load balancing distribution type to be used by the Load Balancer.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(loadbalancers.LoadDistributionDefault),
		},
	}
}
