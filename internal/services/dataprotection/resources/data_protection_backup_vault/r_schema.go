// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataProtectionBackupVaultSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
				"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"datastore_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingStoreTypes(), false),
		},

		"redundancy": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForStorageSettingTypes(), false),
		},

		"cross_region_restore_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retention_duration_in_days": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      14,
			ValidateFunc: validation.FloatBetween(14, 180),
		},

		"soft_delete": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      backupvaultresources.SoftDeleteStateOn,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForSoftDeleteState(), false),
		},

		"immutability": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      backupvaultresources.ImmutabilityStateDisabled,
			ValidateFunc: validation.StringInSlice(backupvaultresources.PossibleValuesForImmutabilityState(), false),
		},

		"encryption_settings": {
			Type:       pluginsdk.TypeList,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Optional:   true,
			// NOTE: O+C - `encryption_settings` is populated when the properties are configured via `azurerm_data_protection_backup_vault_customer_managed_key` resource
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"infrastructure_encryption_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}
