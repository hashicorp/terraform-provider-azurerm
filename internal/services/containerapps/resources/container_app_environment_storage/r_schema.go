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
			ConflictsWith: []string{"nfs_server_url"},
			Description:   "The Azure Storage Account in which the Share to be used is located.",
		},

		"access_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"account_name"},
			AtLeastOneOf:  []string{"access_key", "access_key_wo", "nfs_server_url"},
			ConflictsWith: []string{"access_key_wo", "nfs_server_url"},
			Description:   "The Storage Account Access Key.",
		},

		"access_key_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"account_name", "access_key_wo_version"},
			AtLeastOneOf:  []string{"access_key", "access_key_wo", "nfs_server_url"},
			ConflictsWith: []string{"access_key", "nfs_server_url"},
			Description:   "The Storage Account Access Key.",
		},

		"access_key_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"access_key_wo"},
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "An integer value used to trigger an update for `access_key_wo`. This property should be incremented when updating `access_key_wo`.",
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
			AtLeastOneOf:  []string{"access_key", "access_key_wo", "nfs_server_url"},
			ConflictsWith: []string{"account_name", "access_key", "access_key_wo"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
