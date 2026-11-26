// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dataProtectionBackupInstanceBlobStorageSchema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "The name which should be used for this Backup Instance Blob Storage.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"vault_id": {
			Description:  "The ID of the Backup Vault within which the Backup Instance Blob Storage should exist.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupvaultresources.ValidateBackupVaultID,
		},

		"storage_account_id": {
			Description:  "The ID of the source Storage Account.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"backup_policy_id": {
			Description:  "The ID of the Backup Policy.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupPolicyID,
		},

		"auto_protection_enabled": {
			Description: "Whether all present and future containers of the source Storage Account should be backed up.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"excluded_container_name_prefixes": {
			Description: "A list of container name prefixes to exclude from auto protection. Containers whose names start with one of these prefixes are not backed up.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z][0-9a-z-]*$`), "only lowercase alphanumeric characters and hyphens are allowed, and the value must not start with a hyphen"),
				),
			},
			RequiredWith:  []string{"auto_protection_enabled"},
			ConflictsWith: []string{"storage_account_container_names"},
		},

		"storage_account_container_names": {
			Description: "The list of the container names of the source Storage Account.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			ConflictsWith: []string{"excluded_container_name_prefixes"},
		},

		"protection_state": {
			Description: "The protection state of the Backup Instance Blob Storage.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
