// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_network_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Network Manager Network Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Description:  "The ID of the Network Manager.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkgroups.ValidateNetworkManagerID,
		},

		"description": {
			Description: "A description of the Network Manager Network Group.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"member_type": {
			Description:  "The member type for the network group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(networkgroups.GroupMemberTypeVirtualNetwork),
			ValidateFunc: validation.StringInSlice(networkgroups.PossibleValuesForGroupMemberType(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
