// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Network Manager Admin Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"admin_rule_collection_id": {
			Description:  "The ID of the Network Manager Admin Rule Collection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: adminrules.ValidateSecurityAdminConfigurationRuleCollectionID,
		},

		"action": {
			Description:  "The action allowed for this Network Manager Admin Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(adminrules.PossibleValuesForSecurityConfigurationRuleAccess(), false),
		},

		"direction": {
			Description:  "Whether the traffic matched against the rule in inbound or outbound.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(adminrules.PossibleValuesForSecurityConfigurationRuleDirection(), false),
		},

		"priority": {
			Description:  "The priority of the rule. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 4096),
		},

		"protocol": {
			Description:  "Which network protocol this Network Manager Admin Rule applies to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(adminrules.PossibleValuesForSecurityConfigurationRuleProtocol(), false),
		},

		"description": {
			Description:  "A description of the Network Manager Admin Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"destination_port_ranges": {
			Description: "A list of strings that specify the destination port ranges. Specify one or more single port number or port ranges such as `1024-65535`. Use `*` to specify any port.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"destination": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefix": {
						Description:  "The address prefix.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefix_type": {
						Description: "The address prefix type. For more information, see [this document](https://learn.microsoft.com/azure/virtual-network-manager/concept-security-admins#source-and-destination-types).",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(adminrules.AddressPrefixTypeIPPrefix),
							string(adminrules.AddressPrefixTypeServiceTag),
						}, false),
					},
				},
			},
		},

		"source_port_ranges": {
			Description: "A list of strings that specify the source port ranges. Specify one or more single port number or port ranges such as `1024-65535`. Use `*` to specify any port.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefix": {
						Description:  "The address prefix.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefix_type": {
						Description: "The address prefix type. For more information, see [this document](https://learn.microsoft.com/azure/virtual-network-manager/concept-security-admins#source-and-destination-types).",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(adminrules.AddressPrefixTypeIPPrefix),
							string(adminrules.AddressPrefixTypeServiceTag),
						}, false),
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
