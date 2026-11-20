// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"math"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"storage_sync_group_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: serverendpointresource.ValidateSyncGroupID,
		},

		"registered_server_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: registeredserverresource.ValidateRegisteredServerID,
		},

		"server_local_path": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cloud_tiering_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
			Default:  false,
		},

		"volume_free_space_percent": {
			Optional:     true,
			Type:         pluginsdk.TypeInt,
			Default:      20,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"tier_files_older_than_days": {
			Optional:     true,
			Type:         pluginsdk.TypeInt,
			ValidateFunc: validation.IntBetween(1, math.MaxInt32),
		},

		"initial_download_policy": {
			ForceNew:     true,
			Optional:     true,
			Type:         pluginsdk.TypeString,
			Default:      serverendpointresource.InitialDownloadPolicyNamespaceThenModifiedFiles,
			ValidateFunc: validation.StringInSlice(serverendpointresource.PossibleValuesForInitialDownloadPolicy(), false),
		},

		"local_cache_mode": {
			Optional:     true,
			Type:         pluginsdk.TypeString,
			Default:      serverendpointresource.LocalCacheModeUpdateLocallyCachedFiles,
			ValidateFunc: validation.StringInSlice(serverendpointresource.PossibleValuesForLocalCacheMode(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
