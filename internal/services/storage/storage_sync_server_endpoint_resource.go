// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SyncServerEndpointResource struct{}

var _ sdk.ResourceWithUpdate = SyncServerEndpointResource{}

func (r SyncServerEndpointResource) ModelObject() interface{} {
	return &StorageSyncServerEndpointResourceSchema{}
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

func (r SyncServerEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return serverendpointresource.ValidateServerEndpointID
}

func (r SyncServerEndpointResource) ResourceType() string {
	return "azurerm_storage_sync_server_endpoint"
}

func (r SyncServerEndpointResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validation.IntBetween(1, 2147483647),
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

func (r SyncServerEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SyncServerEndpointResource) Create() sdk.ResourceFunc {
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

			existing, err := client.ServerEndpointsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := serverendpointresource.ServerEndpointCreateParameters{
				Properties: &serverendpointresource.ServerEndpointCreateParametersProperties{
					InitialDownloadPolicy:  pointer.To(serverendpointresource.InitialDownloadPolicy(config.InitialDownloadPolicy)),
					LocalCacheMode:         pointer.To(serverendpointresource.LocalCacheMode(config.LocalCacheMode)),
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

			pollerType := custompollers.NewStorageSyncServerEndpointPoller(client, id)
			poller := pollers.NewPoller(pollerType, 20*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

			if _, err = client.ServerEndpointsCreate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SyncServerEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.SyncServerEndpointsClient

			schema := StorageSyncServerEndpointResourceSchema{}

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

			if model := resp.Model; model != nil {
				schema.Name = id.ServerEndpointName
				if props := model.Properties; props != nil {
					schema.StorageSyncGroupId = serverendpointresource.NewSyncGroupID(id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName, id.SyncGroupName).ID()
					schema.RegisteredServerId = pointer.From(props.ServerResourceId)
					schema.ServerLocalPath = pointer.From(props.ServerLocalPath)
					schema.VolumeFreeSpacePercent = pointer.From(props.VolumeFreeSpacePercent)
					schema.CloudTieringEnabled = pointer.From(props.CloudTiering) == serverendpointresource.FeatureStatusOn
					schema.InitialDownloadPolicy = string(pointer.From(props.InitialDownloadPolicy))
					schema.LocalCacheMode = string(pointer.From(props.LocalCacheMode))
					if pointer.From(props.TierFilesOlderThanDays) != 0 {
						schema.TierFilesOlderThanDays = pointer.From(props.TierFilesOlderThanDays)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r SyncServerEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.SyncServerEndpointsClient

			id, err := serverendpointresource.ParseServerEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.ServerEndpointsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SyncServerEndpointResource) Update() sdk.ResourceFunc {
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
					LocalCacheMode:         pointer.To(serverendpointresource.LocalCacheMode(config.LocalCacheMode)),
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
