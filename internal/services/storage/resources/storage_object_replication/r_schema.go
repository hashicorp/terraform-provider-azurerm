// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageObjectReplicationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"destination_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"rules": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_container_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.StorageContainerName,
					},

					"destination_container_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.StorageContainerName,
					},

					// Possible values are "" means "OnlyNewObjects", "1601-01-01T00:00:00Z" means "Everything" and timeStamp "2020-10-21T16:00:00Z"
					"copy_blobs_created_after": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "OnlyNewObjects",
						ValidateFunc: validate.ObjectReplicationCopyBlobsCreatedAfter,
					},

					"filter_out_blobs_with_prefix": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"metrics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"source_object_replication_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"destination_object_replication_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
