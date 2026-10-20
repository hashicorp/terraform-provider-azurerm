// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}),

		"location": commonschema.Location(),

		"backup_policy_data_lake_storage_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}),

		"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}),

		"storage_container_names": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 1000,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
