// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataProtectionBackupVaultSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Backup Vault.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
				"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"datastore_type": {
			Description:  "The type of the data store.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingStoreTypes(), false),
		},

		"redundancy": {
			Description:  "The backup storage redundancy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingTypes(), false),
		},

		"alerts_for_all_job_failures_enabled": {
			Description: "Whether to enable Azure Monitor alerts for all job failures.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			// Azure Monitor alerts for job failures are enabled by default for Backup vaults.
			// https://learn.microsoft.com/azure/backup/backup-azure-monitoring-alerts?tabs=backup-vaults#turn-on-azure-monitor-alerts-for-job-failure-scenarios
			Default: true,
		},

		"cross_region_restore_enabled": {
			Description: "Whether to enable cross-region restore for the Backup Vault.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"retention_duration_in_days": {
			Description:  "The soft delete retention duration for this Backup Vault.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      14,
			ValidateFunc: validation.FloatBetween(14, 180),
		},

		"soft_delete": {
			Description:  "The state of soft delete for this Backup Vault.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      backupvaultresources.SoftDeleteStateOn,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForSoftDeleteState(), false),
		},

		"immutability": {
			Description:  "The state of immutability for this Backup Vault.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      backupvaultresources.ImmutabilityStateDisabled,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForImmutabilityState(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["alerts_for_all_job_failures_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// NOTE: O+C Preserve existing Azure alert settings when omitted in 5.x rather than enabling alerts on upgrade.
			Computed: true,
		}
	}

	return s
}
