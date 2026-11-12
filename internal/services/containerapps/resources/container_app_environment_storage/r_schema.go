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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentStorageName,
			Description:  "The name for this Storage.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironmentsstorages.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to which this storage belongs.",
		},

		"account_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  storageValidate.StorageAccountName,
			RequiredWith:  []string{"access_key"},
			ConflictsWith: []string{"nfs_server_url"},
			Description:   "The Azure Storage Account in which the Share to be used is located.",
		},

		"access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"account_name"},
			Description:  "The Storage Account Access Key.",
		},

		"share_name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Azure Storage Share to use.",
		},

		"access_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(managedenvironmentsstorages.PossibleValuesForAccessMode(), false),
			Description:  "The access mode to connect this storage to the Container App. Possible values include `ReadOnly` and `ReadWrite`.",
		},

		"nfs_server_url": {
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
