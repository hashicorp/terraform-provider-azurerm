// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_path

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageDataLakeGen2PathSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"filesystem_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: helpers.ValidateStorageDataLakeGen2FileSystemName,
		},

		"path": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true, // TODO handle rename
		},

		"resource": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"directory"}, false),
		},

		"owner": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"$superuser"}, false)),
		},

		"group": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"$superuser"}, false)),
		},

		"ace": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"scope": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "access",
						ValidateFunc: validation.StringInSlice([]string{"default", "access"}, false),
					},
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"user", "group", "mask", "other"}, false),
					},
					"id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"permissions": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ADLSAccessControlPermissions,
					},
				},
			},
		},
	}
}
