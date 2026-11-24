// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_management_private_link_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Optional:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsUUID,
		},
		"management_group_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: commonids.ValidateManagementGroupID,
		},
		"resource_management_private_link_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: resourcemanagementprivatelink.ValidateResourceManagementPrivateLinkID,
		},
		"public_network_access_enabled": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeBool,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"tenant_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
