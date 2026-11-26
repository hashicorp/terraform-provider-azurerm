// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_postgresql_flexible_server

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Backup Instance for the PostgreSQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&basebackuppolicyresources.BackupVaultId{}, "Backup Vault"),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}, "Backup Policy"),

		"server_id": commonschema.ResourceIDReferenceRequiredForceNew(&servers.FlexibleServerId{}, "MySQL Flexible Server"),

		"backup_solution_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(backupinstanceresources.BackupSolutionTypeLogicalBackup),
			ValidateFunc: validation.StringInSlice(backupinstanceresources.PossibleValuesForBackupSolutionType(), false),
		},
	}

	if !features.SixPointOh() {
		arguments["backup_solution_type"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(backupinstanceresources.PossibleValuesForBackupSolutionType(), false),
		}
	}

	return arguments
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Description: "The protection state of the Backup Instance PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if rd == nil || rd.Id() == "" || !rd.NewValueKnown("backup_solution_type") {
				return nil
			}

			oldVal, newVal := rd.GetChange("backup_solution_type")
			if oldVal.(string) == string(backupinstanceresources.BackupSolutionTypePhysicalBackup) &&
				newVal.(string) == string(backupinstanceresources.BackupSolutionTypeLogicalBackup) {
				if err := rd.ForceNew("backup_solution_type"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
