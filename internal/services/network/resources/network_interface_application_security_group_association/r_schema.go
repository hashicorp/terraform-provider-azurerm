// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_security_group_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func networkInterfaceApplicationSecurityGroupAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_interface_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateNetworkInterfaceID,
		},

		"application_security_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
		},
	}
}
