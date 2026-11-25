// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_storage

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name for this Container App Environment Storage.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentStorageName,
		},

		"container_app_environment_id": {
			Description:  "The ID of the Container App Environment to which this storage belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironmentsstorages.ValidateManagedEnvironmentID,
		},

		"account_name": {
			Description:   "The Azure Storage Account in which the Share to be used is located.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  storageValidate.StorageAccountName,
			RequiredWith:  []string{"access_key"},
			ConflictsWith: []string{"nfs_server_url"},
		},

		"access_key": {
			Description:  "The Storage Account Access Key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"account_name"},
		},

		"share_name": {
			Description: "The name of the Azure Storage Share to use.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"access_mode": {
			Description:  "The access mode to connect this storage to the Container App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(managedenvironmentsstorages.PossibleValuesForAccessMode(), false),
		},

		"nfs_server_url": {
			Description:   "The NFS server to use for the Azure File Share, the format will be `yourstorageaccountname.file.core.windows.net`.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"account_name"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
