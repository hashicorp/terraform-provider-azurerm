// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	mssqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseSqlPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SqlPoolName,
		},

		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"DW100c",
				"DW200c",
				"DW300c",
				"DW400c",
				"DW500c",
				"DW1000c",
				"DW1500c",
				"DW2000c",
				"DW2500c",
				"DW3000c",
				"DW5000c",
				"DW6000c",
				"DW7500c",
				"DW10000c",
				"DW15000c",
				"DW30000c",
			}, false),
		},

		"storage_account_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(synapse.PossibleStorageAccountTypeValues(), false),
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  DefaultCreateMode,
			ValidateFunc: validation.StringInSlice([]string{
				DefaultCreateMode,
				RecoveryCreateMode,
				PointInTimeRestoreCreateMode,
			}, false),
		},

		"collation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C The default of this is configurable by the user, so this should remain
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: mssqlValidate.DatabaseCollation(),
		},

		"recovery_database_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validate.SqlPoolID,
				validate.SqlPoolRecoverableDatabaseID,
				commonids.ValidateSqlDatabaseID,
				mssqlValidate.RecoverableDatabaseID,
			),
			ConflictsWith: []string{"restore"},
		},

		"restore": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"recovery_database_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"point_in_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},

					"source_database_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							validate.SqlPoolID,
							commonids.ValidateSqlDatabaseID,
						),
					},
				},
			},
		},

		"data_encrypted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"geo_backup_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}
