// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("/\\?%"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"network_security_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 140),
		},

		"protocol": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleProtocol(), false),
		},

		"source_port_range": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"source_port_ranges"},
		},

		"source_port_ranges": {
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"source_port_range"},
		},

		"destination_port_range": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"destination_port_ranges"},
		},

		"destination_port_ranges": {
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"destination_port_range"},
		},

		"source_address_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
		},

		"source_address_prefixes": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Elem:         &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
		},

		"destination_address_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
		},

		"destination_address_prefixes": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Elem:         &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:          pluginsdk.HashString,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
		},

		// lintignore:S018
		"source_application_security_group_ids": {
			Type:         pluginsdk.TypeSet,
			MaxItems:     10,
			Optional:     true,
			ExactlyOneOf: []string{"source_address_prefix", "source_address_prefixes", "source_application_security_group_ids"},
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			Set: pluginsdk.HashStringInsensitively,
		},

		// lintignore:S018
		"destination_application_security_group_ids": {
			Type:         pluginsdk.TypeSet,
			MaxItems:     10,
			Optional:     true,
			ExactlyOneOf: []string{"destination_address_prefix", "destination_address_prefixes", "destination_application_security_group_ids"},
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			Set: pluginsdk.HashStringInsensitively,
		},

		"access": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleAccess(), false),
		},

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 4096),
		},

		"direction": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securityrules.PossibleValuesForSecurityRuleDirection(), false),
		},
	}
}
