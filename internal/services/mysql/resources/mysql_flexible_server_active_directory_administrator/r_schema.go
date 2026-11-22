// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azureadadministrators.ValidateFlexibleServerID,
		},

		"identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
