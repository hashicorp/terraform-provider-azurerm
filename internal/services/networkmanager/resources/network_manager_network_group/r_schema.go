// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkgroups.ValidateNetworkManagerID,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"member_type": {
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
