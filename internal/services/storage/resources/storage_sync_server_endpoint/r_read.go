// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storageClient.From(metadata.Client).SyncServerEndpointsClient

			id, err := serverendpointresource.ParseServerEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ServerEndpointsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

type StorageSyncServerEndpointResourceSchema struct {
	Name                   string `tfschema:"name"`
	StorageSyncGroupId     string `tfschema:"storage_sync_group_id"`
	RegisteredServerId     string `tfschema:"registered_server_id"`
	ServerLocalPath        string `tfschema:"server_local_path"`
	CloudTieringEnabled    bool   `tfschema:"cloud_tiering_enabled"`
	VolumeFreeSpacePercent int64  `tfschema:"volume_free_space_percent"`
	TierFilesOlderThanDays int64  `tfschema:"tier_files_older_than_days"`
	InitialDownloadPolicy  string `tfschema:"initial_download_policy"`
	LocalCacheMode         string `tfschema:"local_cache_mode"`
}
