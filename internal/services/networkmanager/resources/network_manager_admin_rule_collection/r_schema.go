// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule_collection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrulecollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Network Manager Admin Rule Collection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"security_admin_configuration_id": {
			Description:  "The ID of the Network Manager Security Admin Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: adminrulecollections.ValidateSecurityAdminConfigurationID,
		},

		"network_group_ids": {
			Description: "A list of Network Group IDs which this Network Manager Admin Rule Collection applies to.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: networkgroups.ValidateNetworkGroupID,
			},
		},

		"description": {
			Description: "A description of the Network Manager Admin Rule Collection.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
