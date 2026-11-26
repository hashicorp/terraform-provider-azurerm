// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

		"storage_account_container_names": {
			Description: "The list of the container names of the source Storage Account.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"protection_state": {
			Description: "The protection state of the Backup Instance Blob Storage.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
