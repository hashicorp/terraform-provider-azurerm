// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.SyncServerEndpointsClient

			id, err := serverendpointresource.ParseServerEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config StorageSyncServerEndpointResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := serverendpointresource.ServerEndpointUpdateParameters{
				Properties: &serverendpointresource.ServerEndpointUpdateProperties{
					LocalCacheMode:         pointer.ToEnum[serverendpointresource.LocalCacheMode](config.LocalCacheMode),
					VolumeFreeSpacePercent: pointer.To(config.VolumeFreeSpacePercent),
				},
			}

			cloudTieringEnabled := serverendpointresource.FeatureStatusOff
			if config.CloudTieringEnabled {
				cloudTieringEnabled = serverendpointresource.FeatureStatusOn
			}
			payload.Properties.CloudTiering = pointer.To(cloudTieringEnabled)

			if config.TierFilesOlderThanDays != 0 {
				payload.Properties.TierFilesOlderThanDays = pointer.To(config.TierFilesOlderThanDays)
			}

			pollerType := custompollers.NewStorageSyncServerEndpointPoller(client, *id)
			poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

			if _, err = client.ServerEndpointsUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}
