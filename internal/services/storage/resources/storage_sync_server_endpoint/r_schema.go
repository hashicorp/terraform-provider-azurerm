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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"storage_sync_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: serverendpointresource.ValidateSyncGroupID,
		},

		"registered_server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registeredserverresource.ValidateRegisteredServerID,
		},

		"server_local_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cloud_tiering_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"volume_free_space_percent": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      20,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"tier_files_older_than_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, math.MaxInt32),
		},

		"initial_download_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      serverendpointresource.InitialDownloadPolicyNamespaceThenModifiedFiles,
			ValidateFunc: validation.StringInSlice(serverendpointresource.PossibleValuesForInitialDownloadPolicy(), false),
		},

		"local_cache_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      serverendpointresource.LocalCacheModeUpdateLocallyCachedFiles,
			ValidateFunc: validation.StringInSlice(serverendpointresource.PossibleValuesForLocalCacheMode(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
