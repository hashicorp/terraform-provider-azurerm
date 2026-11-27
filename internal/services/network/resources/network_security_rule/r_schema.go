// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_security_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func networkSecurityRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the security rule. This needs to be unique across all Rules in the Network Security Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("/\\?%"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"network_security_group_name": {
			Description: "The name of the Network Security Group that we want to attach the rule to.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"description": {
			Description:  "A description for this rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 140),
		},

		"protocol": {
			Description:  "Network protocol this rule applies to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleProtocol(), false),
		},

		"source_port_range": {
			Description:   "Source Port or Range. Integer or range between `0` and `65535` or `*` to match any.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"source_port_ranges"},
		},

		"source_port_ranges": {
			Description:   "List of source ports or port ranges.",
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"source_port_range"},
			Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"destination_port_range": {
			Description:   "Destination Port or Range. Integer or range between `0` and `65535` or `*` to match any.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"destination_port_ranges"},
		},

		"destination_port_ranges": {
			Description:   "List of destination ports or port ranges.",
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"destination_port_range"},
			Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"source_address_prefix": {
			Description:  "CIDR or source IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
		},

		"source_address_prefixes": {
			Description:  "List of source address prefixes. Tags may not be used.",
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
			Elem:         &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"destination_address_prefix": {
			Description:  "CIDR or destination IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. Besides, it also supports all available Service Tags like `Sql.WestEurope`, `Storage.EastUS`, etc. You can list the available service tags with the CLI: `az network list-service-tags --location westcentralus`. For further information see [Azure CLI - az network list-service-tags](https://learn.microsoft.com/cli/azure/network?view=azure-cli-latest#az-network-list-service-tags).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
		},

		"destination_address_prefixes": {
			Description:  "List of destination address prefixes. Tags may not be used.",
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
			Elem:         &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		// lintignore:S018
		"source_application_security_group_ids": {
			Description:  "A List of source Application Security Group IDs.",
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			MaxItems:     10,
			Set:          pluginsdk.HashStringInsensitively,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},

		// lintignore:S018
		"destination_application_security_group_ids": {
			Description:  "A List of destination Application Security Group IDs.",
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			MaxItems:     10,
			Set:          pluginsdk.HashStringInsensitively,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},

		"access": {
			Description:  "Whether network traffic is allowed or denied.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleAccess(), false),
		},

		"priority": {
			Description:  "The priority of the rule. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 4096),
		},

		"direction": {
			Description:  "The direction specifies if rule will be evaluated on incoming or outgoing traffic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleDirection(), false),
		},
	}
}
