// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_data_lake_gen2

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataShareDatasetDataLakeGen2Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Data Share Data Lake Gen2 Dataset.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataSetName(),
		},

		"share_id": {
			Description:  "The resource ID of the Data Share where this Data Share Data Lake Gen2 Dataset should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: share.ValidateShareID,
		},

		"storage_account_id": {
			Description:  "The resource ID of the Storage Account of the Data Lake file system to be shared with the receiver.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"file_system_name": {
			Description:  "The name of the Data Lake file system to be shared with the receiver.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"file_path": {
			Description:   "The path of the file in the Data Lake file system to be shared with the receiver.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"folder_path"},
		},

		"folder_path": {
			Description:   "The folder path in the Data Lake file system to be shared with the receiver.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"file_path"},
		},

		"display_name": {
			Description: "The name of the Data Share Dataset.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
