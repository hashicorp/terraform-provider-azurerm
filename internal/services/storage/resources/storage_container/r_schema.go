// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func storageContainerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageContainerName,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"container_access_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "private",
			ValidateFunc: validation.StringInSlice([]string{
				string(containers.Blob),
				string(containers.Container),
				"private",
			}, false),
		},

		"default_encryption_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because needed because a dummy value is returned when unspecified
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageEncryptionScopeName,
		},

		"encryption_scope_override_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      true, // defaulting to false would be preferable here, but the API defaults this to true when unspecified
			RequiredWith: []string{"default_encryption_scope"},
		},

		"metadata": helpers.MetaDataComputedSchema(),

		"has_immutability_policy": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"has_legal_hold": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
