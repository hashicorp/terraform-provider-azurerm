// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := storageClient.From(metadata.Client).StorageDiscoveryWorkspacesClient

			id, err := storagediscoveryworkspaces.ParseProviderStorageDiscoveryWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageDiscoveryWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := storagediscoveryworkspaces.StorageDiscoveryWorkspaceUpdate{
				Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspacePropertiesUpdate{},
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("sku") {
				payload.Properties.Sku = pointer.ToEnum[storagediscoveryworkspaces.StorageDiscoverySku](model.Sku)
			}

			if metadata.ResourceData.HasChange("scope") {
				payload.Properties.Scopes = pointer.To(expandStorageDiscoveryScopes(model.Scope))
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
