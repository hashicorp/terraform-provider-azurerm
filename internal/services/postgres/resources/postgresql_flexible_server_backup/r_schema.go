// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_backup

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/backupautomaticandondemands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this PostgreSQL Flexible Server Backup.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerBackupName,
		},

		"server_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupautomaticandondemands.FlexibleServerId{}, "PostgreSQL Flexible Server"),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"completed_time": {
			Description: "The Time (ISO8601 format) at which the backup was completed.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
