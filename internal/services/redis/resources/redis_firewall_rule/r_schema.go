// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_firewall_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func redisFirewallRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallRuleName,
		},

		"redis_cache_name": {
			Description: "The name of the Redis Cache.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"start_ip": {
			Description: "The lowest IP address included in the range.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.IsIPAddress,
				validation.StringIsNotEmpty,
			),
		},

		"end_ip": {
			Description: "The highest IP address included in the range.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.IsIPAddress,
				validation.StringIsNotEmpty,
			),
		},
	}
}
