// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_data_lake_storage

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Backup Instance Data Lake Storage.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}, "Data Protection Backup Vault"),

		"location": commonschema.Location(),

		"backup_policy_data_lake_storage_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}, "Backup Policy"),

		"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}, "Storage Account"),

		"storage_container_names": {
			Description: "The names of the Storage Containers of the source Storage Account.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			MaxItems:    1000,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 63),
					validation.StringMatch(regexp.MustCompile(`^[0-9a-z][0-9a-z-]*$`), "only lowercase alphanumeric characters and hyphens are allowed, and the value must not start with a hyphen"),
				),
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Description: "The protection state of the Backup Instance Data Lake Storage.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
