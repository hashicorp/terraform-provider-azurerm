// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_discovery_workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageDiscoveryWorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model StorageDiscoveryWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := storagediscoveryworkspaces.NewProviderStorageDiscoveryWorkspaceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := storagediscoveryworkspaces.StorageDiscoveryWorkspace{
				Location: location.Normalize(model.Location),
				Properties: &storagediscoveryworkspaces.StorageDiscoveryWorkspaceProperties{
					WorkspaceRoots: model.WorkspaceRoots,
					Scopes:         expandStorageDiscoveryScopes(model.Scope),
					Sku:            pointer.ToEnum[storagediscoveryworkspaces.StorageDiscoverySku](model.Sku),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func expandStorageDiscoveryScopes(input []StorageDiscoveryScopeModel) []storagediscoveryworkspaces.StorageDiscoveryScope {
	result := make([]storagediscoveryworkspaces.StorageDiscoveryScope, 0)

	for _, scope := range input {
		apiScope := storagediscoveryworkspaces.StorageDiscoveryScope{
			DisplayName:   scope.DisplayName,
			ResourceTypes: pointer.From(pointer.ToEnumSlice[storagediscoveryworkspaces.StorageDiscoveryResourceType](scope.ResourceTypes)),
		}

		if len(scope.TagKeysOnly) > 0 {
			apiScope.TagKeysOnly = &scope.TagKeysOnly
		}

		if len(scope.Tags) > 0 {
			apiScope.Tags = &scope.Tags
		}

		result = append(result, apiScope)
	}

	return result
}
