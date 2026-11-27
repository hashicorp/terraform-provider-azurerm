// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet_network_security_group_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func subnetNetworkSecurityGroupAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subnet_id": {
			Description:  "The ID of the Subnet.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"network_security_group_id": {
			Description:  "The ID of the Network Security Group which should be associated with the Subnet.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
		},
	}
}
