// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *serverendpointresource.ServerEndpointId, model *serverendpointresource.ServerEndpoint) error {
	schema := StorageSyncServerEndpointResourceSchema{
		Name:               id.ServerEndpointName,
		StorageSyncGroupId: serverendpointresource.NewSyncGroupID(id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName, id.SyncGroupName).ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			schema.RegisteredServerId = pointer.From(props.ServerResourceId)
			schema.ServerLocalPath = pointer.From(props.ServerLocalPath)
			schema.VolumeFreeSpacePercent = pointer.From(props.VolumeFreeSpacePercent)
			schema.CloudTieringEnabled = pointer.From(props.CloudTiering) == serverendpointresource.FeatureStatusOn
			schema.InitialDownloadPolicy = pointer.FromEnum(props.InitialDownloadPolicy)
			schema.LocalCacheMode = pointer.FromEnum(props.LocalCacheMode)
			if pointer.From(props.TierFilesOlderThanDays) != 0 {
				schema.TierFilesOlderThanDays = pointer.From(props.TierFilesOlderThanDays)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&schema)
}
