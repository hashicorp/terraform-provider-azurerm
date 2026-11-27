// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.SyncServerEndpointsClient

			var config StorageSyncServerEndpointResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			storageSyncGroupId, err := serverendpointresource.ParseSyncGroupID(config.StorageSyncGroupId)
			if err != nil {
				return err
			}

			id := serverendpointresource.NewServerEndpointID(subscriptionId, storageSyncGroupId.ResourceGroupName, storageSyncGroupId.StorageSyncServiceName, storageSyncGroupId.SyncGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ServerEndpointsGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := serverendpointresource.ServerEndpointCreateParameters{
				Properties: &serverendpointresource.ServerEndpointCreateParametersProperties{
					InitialDownloadPolicy:  pointer.ToEnum[serverendpointresource.InitialDownloadPolicy](config.InitialDownloadPolicy),
					LocalCacheMode:         pointer.ToEnum[serverendpointresource.LocalCacheMode](config.LocalCacheMode),
					ServerLocalPath:        pointer.To(config.ServerLocalPath),
					ServerResourceId:       pointer.To(config.RegisteredServerId),
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

			// TODO: confirm whether this poller is still required, the go-azure-sdk LRO poller has changed since
			pollerType := custompollers.NewStorageSyncServerEndpointPoller(client, id)
			poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

			if _, err = client.ServerEndpointsCreate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}
