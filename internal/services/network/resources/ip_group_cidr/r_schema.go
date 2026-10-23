// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group_cidr

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ipGroupCidrSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ip_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ipgroups.ValidateIPGroupID,
		},
		"cidr": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
